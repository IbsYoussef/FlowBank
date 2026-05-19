from enum import Enum
from pydantic import BaseModel, Field
from datetime import datetime, timezone
from typing import Optional, List

class TransactionType(str, Enum):
    DEBIT = "debit"
    CREDIT = "credit"

class RiskScore(str, Enum):
    LOW = "low"
    MEDIUM = "medium"
    HIGH = "high"

class FraudStatus(str, Enum):
    CLEAN = "clean"
    SUSPICIOUS = "suspicious"
    FLAGGED = "flagged"

class Transaction(BaseModel):
    id: str
    user_id: str
    amount: int
    type: TransactionType
    description: str
    merchant_name: str
    status: str
    created_at: datetime

class FraudScore(BaseModel):
    transaction_id: str
    risk_score: RiskScore
    status: FraudStatus
    triggered_rules: List[str] = Field(default_factory=list)
    confidence: float = Field(ge=0.0, le=1.0)
    processing_time_ms: Optional[float] = None
    scored_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    model_config = {
        "use_enum_values": True
    }