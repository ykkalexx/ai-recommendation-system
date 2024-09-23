from flask import Flask, request, jsonify
from model import RecommendationModel

app = Flask(__name__)
model = RecommendationModel()

@app.route('/train', methods=['POST'])
def train_model():
    # Soon, we would trigger model training here . For now, ill just load data and train
    model.load_data('create_csv_and add_it_here')  
    model.train()
    return jsonify({"message": "Model trained successfully"}), 200

@app.route('/recommend', methods=['GET'])
def get_recommendations():
    user_id = request.args.get('user_id')
    if not user_id:
        return jsonify({"error": "User ID is required"}), 400
    
    recommendations = model.get_recommendations(user_id)
    return jsonify({"user_id": user_id, "recommendations": recommendations}), 200

if __name__ == "__main__":
    app.run(debug=True, port=5000)