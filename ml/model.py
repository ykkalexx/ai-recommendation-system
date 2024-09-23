import pandas as pd
from surprise import Dataset, Reader, SVD
from surprise.model_selection import train_test_split
from surprise import accuracy
from pymongo import MongoClient

class RecommendationModel:
    def __init__(self, mongo_uri):
        self.model = SVD()
        self.data = None
        self.mongo_client = MongoClient(mongo_uri)

    def load_data(self):
        behaviors = list(self.db.behaviors.find({}, {'_id': 0, 'user_id': 1, 'item_id': 1, 'action': 1}))
        df = pd.DataFrame(behaviors)

        # convert 'action' to numeric rating (e.g 'view' = 1)
        df['rating']= df['action'].map({'view': 1})

    def train(self):
        trainset = self.data.build_full_trainset()
        self.model.fit(trainset)

    def get_recommendations(self, user_id, n=5):
        # Get all items
        all_items = self.data.df['item_id'].unique()
        # Predict ratings for all items
        item_predicitons = [self.model.predict(user_id, item) for item in all_items]
        # sort predictions by estimated rating
        item_predicitons.sort(key=lambda x: x.est, reverse=True)

        # return top n recommendations
        return [pred.iid for pred in item_predicitons[:n]]

if __name__ == "__main__":
    model = RecommendationModel()
    model.load_data('create_csv_and add_it_here')  
    model.train()
    
    user_id = "user1"  # Replace with the actual user id once everything is fully setup
    recommendations = model.get_recommendations(user_id)
    print(f"Recommendations for user {user_id}: {recommendations}")