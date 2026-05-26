# FlowBank - Polish & Future Improvements

This document tracks planned improvements and polish items for the FlowBank fraud detection feature. These are post-deployment additions to enhance the demo experience and technical depth.

---

## 🎨 Dashboard Polish

### Pagination

- Add pagination to the Live Transaction Feed and Flagged Transactions panel
- Currently loads latest 10/50 entries - paginating allows browsing full history
- Consider infinite scroll as an alternative for a smoother UX

### Service Status Indicators

- Add a status bar showing connectivity of all services in real time
- Go API ✅ | Fraud Detection ✅ | Kafka ✅ | Database ✅
- Visible at a glance, tells the distributed system story immediately

### Performance

- Reduce the 3 second Kafka consumer delay currently needed to avoid race conditions
- Implement a retry mechanism instead of a fixed delay for faster real-time updates
- Dashboard feed should feel snappier and more live

---

## 🎲 Data Variety

### Transaction Generator (Producer)

- Current producer generates very similar transactions - limited merchant variety and amount ranges
- Add more diverse merchant categories (e.g. travel, groceries, entertainment, utilities)
- Add more realistic amount distributions (small frequent purchases vs occasional large ones)
- Add more test users beyond Alice and Bob to create richer patterns

### Fraud Detection Scenarios

- Currently almost all flagged transactions trigger only `high_frequency` rule
- Seed specific scenarios to trigger all three rules visibly:
  - **Duplicate transactions**: same user, same amount, within 60 seconds
  - **High value**: transactions above $10,000
  - **High frequency**: more than 5 transactions in 60 seconds
- Add more green (low risk) transactions so the contrast with flagged ones is clearer
- Consider adding a demo mode that deliberately triggers all fraud rules for visual impact

---

## ⚙️ CI/CD Pipeline

### GitHub Actions

- Set up workflow that triggers on push to `main`
- Pipeline steps:
  - Run Go tests
  - Build Docker images
  - Push to ECR (AWS Container Registry)
  - Deploy to Elastic Beanstalk
- Add a build status badge to the README
- Talking point: "Every push to main automatically ships to production"

---

## 🏗️ Infrastructure

### ECS Migration (longer term)

- Migrate from Elastic Beanstalk to ECS for better container orchestration
- Stronger AWS talking point for interviews
- Migration story itself is a portfolio talking point
- Do this after the initial Elastic Beanstalk deployment is stable

### Monitoring

- Add Prometheus metrics endpoint to fraud detection service
- Basic latency and throughput visibility
- Talking point: "Sub-100ms p95 fraud scoring latency under load"

---

## Priority Order

| Item                      | Priority | Effort | Impact |
| ------------------------- | -------- | ------ | ------ |
| CI/CD pipeline            | High     | Low    | High   |
| Data variety              | High     | Medium | High   |
| Reduce consumer delay     | Medium   | Low    | Medium |
| Pagination                | Medium   | Low    | Medium |
| Service status indicators | Medium   | Low    | High   |
| ECS migration             | Low      | High   | High   |
| Monitoring                | Low      | Medium | Medium |
