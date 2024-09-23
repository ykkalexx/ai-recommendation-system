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
        self.db = self.mongo_client.recommandation-system

    def load_data(self):
        behaviors = list(self.db.behaviors.find({}, {'_id': 0, 'user_id': 1, 'item_id': 1, 'action': 1}))
        df = pd.DataFrame(behaviors)

        # convert 'action' to numeric rating (e.g 'view' = 1)
        df['rating'] = df['action'].map({
            'view': 1,
            'like': 2,
            'comment': 3,
            'share': 4,
            'purchase': 5
        })

        # assuming binnary interaction for now
        reader = Reader(rating_scale=(1,1)) 
        self.data = Dataset.load_from_df(df[['user_id', 'item_id', 'rating']], reader)
        
    def train(self):
        trainset = self.data.build_full_trainset()
        self.model.fit(trainset)

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
        
if __name__ == "__main__":
    model = RecommendationModel('mongodb+srv://alex:alex@cluster0.t5mgz.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0')
    model.load_data()  
    model.train()
    