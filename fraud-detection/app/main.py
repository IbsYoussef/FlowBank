from fastapi import FastAPI

app = FastAPI()

@app.get("/")
async def root():
    return {"message": "root pathway working"}

@app.get("/health")
async def health():
    return {"health": "healthy"}