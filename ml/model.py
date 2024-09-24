import pandas as pd
from surprise import Dataset, Reader, SVD
from surprise.model_selection import train_test_split
from surprise import accuracy
from pymongo import MongoClient
from dotenv import load_dotenv
import os

class RecommendationModel:
    def __init__(self, mongo_uri):
        self.model = SVD()
        self.data = None
        self.mongo_client = MongoClient(mongo_uri)
        self.db = self.mongo_client.recommendationDB

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
        # Get all items
        all_items = self.data.df['item_id'].unique()
        # get the items the user has already interacted with
        user_items = set(self.db.behaviors.distinct('item_id', {'user_id': user_id}))
        # items that the user hasn't interacted with
        candidate_items = list(all_items - user_items)
        # get the predictions for the candidate items
        items_predictions = [self.model.predict(user_id, item) for item in candidate_items]
        # sort them by estimated rating
        items_predictions.sort(key=lambda x: x.est, reverse=True)

        # get the top n recommendations
        return [pred.iid for pred in items_predictions[:n]]
        print("Recommendations: ", recommendations)
        
if __name__ == "__main__":
    load_dotenv()
    mongo_uri = os.getenv('MONGO_URI')
    model = RecommendationModel(mongo_uri)
    model.load_data()  
    model.train()
    