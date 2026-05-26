import logging
import time
from app.models import Transaction, FraudScore, RiskScore, FraudStatus

logger = logging.getLogger(__name__)

class FraudDetector:
    def __init__(self, db):
        self.db = db
        logger.info("Fraud detector initialised")
    
    async def _check_duplicate(self, transaction: Transaction) -> bool:
        result = await self.db.check_duplicate_transaction(
            user_id=transaction.user_id,
            amount=transaction.amount,
            transaction_id=transaction.transaction_id
        )
        return result
    
    async def _check_high_frequency(self, transaction: Transaction) -> bool:
        result = await self.db.get_recent_transaction_count(
            user_id=transaction.user_id,
            seconds=60
        )
        return result > 5
    
    def _calculate_risk(self, triggered_rules: list) -> tuple:
        rule_count = len(triggered_rules)

        if rule_count == 0:
            return RiskScore.LOW, FraudStatus.CLEAN
        elif rule_count == 1:
            return RiskScore.MEDIUM, FraudStatus.SUSPICIOUS
        else:
            return RiskScore.HIGH, FraudStatus.FLAGGED
    
    def _check_high_value(self, transaction: Transaction) -> bool:
        return transaction.amount > 1_000_000
    
    def _calculate_confidence(self, triggered_rules: list) -> float:
        rule_count = len(triggered_rules)

        if rule_count == 0:
            return 0.95
        elif rule_count == 1:
            return 0.75
        elif rule_count == 2:
            return 0.90
        else:
            return 0.98
    
    async def score_transaction(self, transaction: Transaction) -> FraudScore:
        start_time = time.time()
        triggered_rules = []

        if await self._check_duplicate(transaction):
            triggered_rules.append("duplicate_transaction")

        if self._check_high_value(transaction):
            triggered_rules.append("high_value")

        if await self._check_high_frequency(transaction):
            triggered_rules.append("high_frequency")

        risk_score, status = self._calculate_risk(triggered_rules)
        processing_time_ms = (time.time() - start_time) * 1000

        return FraudScore(
            transaction_id=transaction.transaction_id,
            risk_score=risk_score,
            status=status,
            triggered_rules=triggered_rules,
            confidence=self._calculate_confidence(triggered_rules),
            processing_time_ms=processing_time_ms
        )