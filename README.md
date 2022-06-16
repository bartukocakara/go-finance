<div style="display:flex">
  <img src="/public/bg_banner.png" alt="Alt text" title="Optional title" width="120">
  <img src="/public/go-logo.png" alt="Alt text" title="Optional title" width="120">
</div>
### Finance App
- By using this application user can track all expenses, income, etc...
- User can create different accounts and create transaction. Can create budgets scheduled income payments


### Packages

- POSTGRES - Database https://github.com/jmoiron/sqlx (https://github.com/lib/pq)
- MIGRATE - Migration https://github.com/golang-migrate/migrate
- MUX - Routing https://github.com/gorilla/mux
- LOGRUS - Logging  https://github.com/sirupsen/logrus


### Technologies

- Backend will be deployed to AWS or Heroku
- Deployment will be done by Docker we will choose AWS, environment setup will be done by Terraform


### Features
* [ ] User create profile email, facebook or google 
* [ ] User create transaction in each account, transaction can be income, expenses or transfer to other account
* User can create categories and assign transactions to different categories
* [ ] User can see total amount of money on each account and trends


### Run project locally with Docker
```
docker-compose build
docker-compose up

```

