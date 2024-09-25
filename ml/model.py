import pandas as pd
from surprise import Dataset, Reader, SVD
from surprise.model_selection import train_test_split
from surprise import accuracy
from pymongo import MongoClient
from dotenv import load_dotenv
import os
import logging

class RecommendationModel:
    def __init__(self, mongo_uri):
        self.model = SVD()
        self.data = None
        self.mongo_client = MongoClient(mongo_uri)
        self.db = self.mongo_client.recommendationDB
        self.load_data()
        self.train()

    def load_data(self):
        behaviors = list(self.db.behaviors.find({}, {'_id': 0, 'userid': 1, 'itemid': 1, 'action': 1}))
        print("Behaviors: ", behaviors)
        df = pd.DataFrame(behaviors)
        print(df)
        # convert 'action' to numeric rating (e.g 'purchase' = 5)
        df['rating'] = df['action'].map({
            'view': 1,
            'like': 2,
            'comment': 3,
            'share': 4,
            'purchase': 5
        })

        # assuming binary interaction for now
        reader = Reader(rating_scale=(1,1)) 
        self.data = Dataset.load_from_df(df[['userid', 'itemid', 'rating']], reader)
        
    def train(self):
        trainset = self.data.build_full_trainset()
        self.model.fit(trainset)
        print("Model trained successfully")

    def get_recommendations(self, user_id, n=5):
        if self.data is None:
            logging.error("Data has not been loaded. Call load_data before calling get_recommendations.")
            return []

        # Get all items
        all_items = self.data.df['itemid'].unique()
        # get the items the user has already interacted with
        user_items = set(self.db.behaviors.distinct('itemid', {'userid': user_id}))
        # items that the user hasn't interacted with
        candidate_items = list(set(all_items) - user_items)

        if not candidate_items:
            logging.info(f"No candidate items found for user {user_id}")
            return []

        # get the predictions for the candidate items
        items_predictions = [self.model.predict(user_id, item) for item in candidate_items]
        # sort them by estimated rating
        items_predictions.sort(key=lambda x: x.est, reverse=True)

        # get the top n recommendations
        recommendations = [pred.iid for pred in items_predictions[:n]]
        logging.info(f"Recommendations for user {user_id}: {recommendations}")
        return recommendations
        
if __name__ == "__main__":
    load_dotenv()
    mongo_uri = os.getenv('MONGO_URI')
    model = RecommendationModel(mongo_uri)

    