import schedule
import time
from model import RecommendationModel
import os
from dotenv import load_dotenv

def train_model():
    load_dotenv()
    mongo_uri = os.getenv('MONGO_URI')
    model = RecommendationModel(mongo_uri)
    model.load_data()
    model.train()
    print("Model training completed successfully")

def run_pipeline():
    schedule.every(1).day.do(train_model)

    while True:
        schedule.run_pending()
        time.sleep(1)

if __name__ == "__main__":
    run_pipeline()