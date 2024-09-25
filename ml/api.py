import logging
from flask import Flask, request, jsonify
from model import RecommendationModel
from dotenv import load_dotenv
import os
from celery import Celery
from flask_cors import CORS

# Configure logging
logging.basicConfig(
    level=logging.DEBUG,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    handlers=[
        logging.FileHandler("app.log"),
        logging.StreamHandler()
    ]
)

load_dotenv()
mongo_uri = os.getenv('MONGO_URI')

app = Flask(__name__)
model = RecommendationModel(mongo_uri)

# Configure Celery to use Redis as the broker
app.config['CELERY_BROKER_URL'] = 'redis://localhost:6379/0'
app.config['CELERY_RESULT_BACKEND'] = 'redis://localhost:6379/0'
#dealing with cors
CORS(app, origins=['http://localhost:8080'])

celery = Celery(app.name, broker=app.config['CELERY_BROKER_URL'])
celery.conf.update(app.config)

@celery.task(bind=True)
def train_model_task(self):
    logging.info("Starting to train model")
    model.load_data()
    logging.info("Data is loaded")
    model.train()
    return {"message": "Model trained successfully"}

@app.route('/train', methods=['POST'])
def train_model():
    train_model_task.delay()
    return jsonify({"message": "Training model started"}), 202

@app.route('/recommend', methods=['GET'])
def get_recommendations():
    user_id = request.args.get('user_id')
    if not user_id:
        return jsonify({"error": "User ID is required"}), 400
    try:
        recommendations = model.get_recommendations(user_id)
        return jsonify({"user_id": user_id, "recommendations": recommendations}), 200
    except Exception as e:
        logging.error(f"Error getting recommendations for user {user_id}: {e}")
        return jsonify({"error": "Failed to get recommendations"}), 500

if __name__ == "__main__":
    app.run(debug=True, port=5000)