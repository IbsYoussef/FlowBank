import logging
import asyncpg

logger = logging.getLogger(__name__)

class FraudDatabase:
    def __init__(self, connection_string: str):
        self.connection_string = connection_string
        self.pool = None
        logger.info("Database initialised")

    async def connect(self):
        self.pool = await asyncpg.create_pool(
            self.connection_string,
            min_size=2,
            max_size=10
        )
        logger.info("Database connection pool created")
    
    async def close(self):
        if self.pool:
            await self.pool.close()
            logger.info("Database connection pool closed")
    
    async def insert_fraud_score(self, fraud_score) -> bool:
        try:
            async with self.pool.acquire() as conn:
                await conn.execute("""
                    INSERT INTO fraud_scores (
                        transaction_id, risk_score, status,
                        triggered_rules, confidence, processing_time_ms, scored_at
                    ) VALUES ($1, $2, $3, $4, $5, $6, $7)
                """,
                    fraud_score.transaction_id,
                    fraud_score.risk_score,
                    fraud_score.status,
                    fraud_score.triggered_rules,
                    fraud_score.confidence,
                    fraud_score.processing_time_ms,
                    fraud_score.scored_at
                )
                logger.info(f"Inserted fraud score for transaction {fraud_score.transaction_id}")
                return True
        except Exception as e:
            logger.error(f"Error inserting fraud score: {e}")
            return False
    
    async def check_duplicate_transaction(
            self,
            user_id: str,
            amount: int,
            transaction_id: str,
            seconds: int = 60
    ) -> bool:
        try:
            async with self.pool.acquire() as conn:
                result = await conn.fetchval(f"""
                    SELECT COUNT(*) FROM transactions
                    WHERE user_id = $1
                    AND amount = $2
                    AND id != $3
                    AND created_at > NOW() - INTERVAL '{seconds} seconds'
                                             """, user_id, amount, transaction_id)
                return result > 0
        except Exception as e:
            logger.error(f"Error checking duplicate transaction: {e}")
            return False
    
    async def get_recent_transaction_count(
            self,
            user_id: str,
            seconds: int = 60
    ) -> int:
        try:
            async with self.pool.acquire() as conn:
                result = await conn.fetchval(f"""
                    SELECT COUNT(*) FROM transactions
                    WHERE user_id = $1
                    AND created_at > NOW() - INTERVAL '{seconds} seconds'
""", user_id)
                return result
        except Exception as e:
            logger.error(f"Error getting recent transaction count: {e}")
            return 0
            