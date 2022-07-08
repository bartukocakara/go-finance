<div style="display:flex">
  <img src="/public/bg_banner.png" alt="Alt text" title="Finance App" width="120">
  <img src="/public/go-logo.png" alt="Alt text" title="Golang" width="120">
  <img src="/public/postgres-logo.png" alt="Alt text" title="PostgreSQL" width="120">
  <img src="/public/docker.png" alt="Alt text" title="Docker" width="120">
</div>

### Semantic Versioning 2.0.0 https://semver.org/lang/tr/

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
* [ ] Multiple currency options

### Configs
* [ ] Localization
* [ ] DB Access
* [ ] Timezone

### Run project locally with Docker
```
make build-dev
make up-dev

```

### Endpoints
#### BASE URL = http://localhost:8088
#### PREFIX = api
#### VERSION = v1
### FULL URL = http://localhost:8088/api/v1/
| Endpoints  | Description |  Methods | Params | Header | Allowed Roles |
| :------:|  :-----------:| :-----------:| :-----------:| :-----------:| :-----------:|
| /users   | Create User  | POST | email, deviceID, password| - | - |
| /users   | List Users | GET | - | Bearer {Token} | 'admin' |
| /login   | Login User  | POST | email, deviceID, password | - | - |
| /users/{userID}/roles   | Grant Role to User | POST | role(user, admin) | Bearer {Token} | 'admin' |
| /users/{userID}/roles   | Revoke Role from User | DELETE | role(user, admin) | Bearer {Token} | 'admin' |
| /users/{userID}/roles   | List Users Role | GET | - | Bearer {Token} | 'admin' |
| /users/{userID}/roles   | Update Users Role  | UPDATE | role(user, admin) | Bearer {Token} | 'admin' |

