<div style="display:flex">
  <img src="/public/bg_banner.png" alt="Alt text" title="Finance App" width="120">
  <img src="/public/go-logo.png" alt="Alt text" title="Golang" width="120">
  <img src="/public/postgres-logo.png" alt="Alt text" title="PostgreSQL" width="120">
  <img src="/public/docker.png" alt="Alt text" title="Docker" width="120">
</div>

# Finance App
| No. | Topic                                                                   |
| --- | ----------------------------------------------------------------------- |
| 1   | [**Purpose**](#Purpose)                             |
| 2   | [**Packages**](#Packages)                               |
| 3   | [**Technologies**](#Technologies)                                                 |
| 4   | [**Features**](#Features)                                                 |
| 5   | [**Configs**](#Configs)                                                 |
| 6   | [**Endpoints**](#Endpoints)                                                 |

### Semantic Versioning 2.0.0 https://semver.org/lang/tr/

### Purpose
- By using this application user can track all expenses, income, etc...
- User can create different accounts and create transaction. Can create budgets scheduled income payments


### Packages

- POSTGRES - Database https://github.com/jmoiron/sqlx (https://github.com/lib/pq)
- MIGRATE - Migration https://github.com/golang-migrate/migrate
- MUX - Routing https://github.com/gorilla/mux
- LOGRUS - Logging  https://github.com/sirupsen/logrus
- CACHE - Caching github.com/bluele/gcache

### Server

- Backend will be deployed to AWS or Heroku

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
| /users   | List Users | GET | - | Bearer {Token} | admin |
| /users/{UserID}   | Get User | GET | - | Bearer {Token} | admin |
| /users/{UserID}   | Update User | PATCH | email, password | Bearer {Token} | admin |
| /users/{UserID}   | Delete User | DELETE | - | Bearer {Token} | admin |
| /login   | Login User  | POST | email, deviceID, password | - | - |
| /users/{UserID}/roles   | Grant Role to User | POST | role(user, admin) | Bearer {Token} | admin |
| /users/{UserID}/roles   | Revoke Role from User | DELETE | role(user, admin) | Bearer {Token} | admin |
| /users/{UserID}/roles   | List Users Role | GET | - | Bearer {Token} | admin |
| /users/{UserID}/roles   | Update Users Role  | PATCH | role(user, admin) | Bearer {Token} | admin |
| /users/{UserID}/merchants   | Create Merchant  | POST | name | Bearer {Token} | admin |
| /users/{UserID}/merchants/{MerchantID}   | Get Merchant  | GET | - | Bearer {Token} | admin |
| /users/{UserID}/merchants/{MerchantID}   | Update Merchant  | PATCH | name | Bearer {Token} | admin |
| /users/{UserID}/merchants   | List Merchant  | GET | - | Bearer {Token} | admin |
| /users/{UserID}/merchants/{MerchantID}   | Delete Merchant  | DELETE | - | Bearer {Token} | admin,user |
| /users/{UserID}/categories   | Create Category  | POST |  name | Bearer {Token} | admin,user |
| /users/{UserID}/categories/{CategoryID}   | Update Category  | PATCH | name | Bearer {Token} | admin,user |
| /users/{UserID}/categories/{CategoryID}   | Get Category  | GET | - | Bearer {Token} | admin,user |
| /users/{UserID}/categories   | List Category  | GET | - | Bearer {Token} | admin,user |
| /users/{UserID}/categories/{CategoryID}   | Delete Category  | DELETE | - | Bearer {Token} | admin,user |
| /users/{UserID}/accounts   | Create Account  | POST |  account_name,type,start_balance,currency | Bearer {Token} | admin,user |
| /users/{UserID}/accounts/{AccountID}   | Update Account  | PATCH | account_name,type,start_balance,currency | Bearer {Token} | admin,user |
| /users/{UserID}/accounts/{AccountID}   | Get Account  | GET | - | Bearer {Token} | admin,user |
| /users/{UserID}/accounts   | List Category  | GET | - | Bearer {Token} | admin,user |
| /users/{UserID}/accounts/{AccountID}   | Delete Account  | DELETE | - | Bearer {Token} | admin,user |
| /users/{UserID}/transactions   | Create Account  | POST |  account_id,category_id,type,amount,currency,notes | Bearer {Token} | admin,user |
| /transactions   | List All Transactions  | GET | (Q)from,to | Bearer {Token} | admin,user |
| /transactions/{TransactionID}   | Show Transaction  | GET | - | Bearer {Token} | admin,user |
| /accounts/{AccountID}/transactions   | Get Transaction By Account ID  | GET | (Q)from,to | Bearer {Token} | admin,user |
| /categories/{CategoryID}/transactions   | List Transaction By Category ID  | GET | (Q)from,to | Bearer {Token} | admin,user |
| /users/{UserID}/transactions  | List Transaction By User ID  | GET | (Q)from,to | Bearer {Token} | admin,user |
| /users/{UserID}/transactions/{TransactionID}   | Get Transaction By ID  | GET | - | Bearer {Token} | admin,user |
| /users/{UserID}/transactions/{TransactionID}   | Update Transaction By ID  | PATCH | account_id,category_id,type,date,currency,amount | Bearer {Token} | admin,user |
| /users/{UserID}/transactions/{TransactionID}   | Delete Transaction By ID  | DELETE | - | Bearer {Token} | admin,user |