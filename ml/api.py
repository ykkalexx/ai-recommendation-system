from flask import Flask, request, jsonify
from model import RecommendationModel
from dotenv import load_dotenv
import os
from celery import Celery

load_dotenv()
mongo_uri = os.getenv('MONGO_URI')

app = Flask(__name__)
model = RecommendationModel(mongo_uri)

app.config['CELERY_BROKER_URL'] = 'redis://localhost:6379/0'
app.config['CELERY_RESULT_BACKEND'] = 'redis://localhost:6379/0'
celery = Celery(app.name, broker=app.config['CELERY_BROKER_URL'])
celery.conf.update(app.config)

@celery.task(bind=true)
def train_model_task(self):
    model.load_data()
    model.train()
    return {"message": "Model trained successfully"}

@app.route('/train', methods=['POST'])
def train_model():
    model.load_data()
    model.train()
    return {"message": "Model trained successfully"}, 200

@app.route('/recommend', methods=['GET'])
def get_recommendations():
    user_id = request.args.get('user_id')
    if not user_id:
        return jsonify({"error": "User ID is required"}), 400
    
    recommendations = model.get_recommendations(user_id)
    return jsonify({"user_id": user_id, "recommendations": recommendations}), 200

if __name__ == "__main__":
    app.run(debug=True, port=5000)