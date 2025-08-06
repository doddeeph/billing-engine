# Billing Engine

We are building a billing system for our Loan Engine. Basically the job of a billing engine is to provide the:
- Loan schedule for a given loan (when am i supposed to pay how much)
- Outstanding Amount for a given loan
- Status of weather the customer is Delinquent or not

We offer loans to our customers a 50 week loan for Rp 5,000,000/- , and a flat interest rate of 10% per annum.

This means that when we create a new loan for the customer (say loan id 100) then it needs to provide the billing schedule for the loan as

```text
W1: 110000
W2: 110000
W3: 110000
…
W50: 110000
```

The Borrower repays the Amount every week. (assume that borrower can only pay the exact amount of payable that week or not pay at all)

We need the ability to track the Outstanding balance of the loan (defined as pending amount the borrower needs to pay at any point) eg. at the beginning of the week it is 5,500,000/- and it decreases as the borrower continues to make repayment, at the end of the loan it should be 0/-

Some customers may miss repayments, If they miss 2 continuous repayments they are delinquent borrowers.

To cover up for missed payments customers will need to make repayments for the remaining amounts. ie if there are 2 pending payments they need to make 2 repayments (of the exact amount).

We need to track the borrower if the borrower is  (any borrower that who’s not paid for last 2 repayments)

We are looking for at least the following methods to be implemented:
- `GetOutstanding`: This returns the current outstanding on a loan, 0 if no outstanding (or closed),
- `IsDelinquent`: If there are more than 2 weeks of Non payment of the loan amount
- `MakePayment`: Make a payment of certain amount on the loan

## Integration Test
- Go to the project root folder
    ```text
    $ pwd
    ../billing-engine
    ```
- running the project integration tests
    ```text
    $ go test -v ./internal/billing/
    ```

## Running the project
- Go to the project root folder
    ```text
    $ pwd
    ../billing-engine
    ```
- Running the project with docker compose
    ```text
    $ docker compose up -d --build
    ```

## REST API
- Create Billing
    Request:
    ```curl
    curl -X POST http://localhost:8080/api/v1/billings \
    -H "Content-Type: application/json" \
    -d '{
        "customerId": 1,
        "loanId": 1001,
        "loanAmount": 5000000,
        "loanInterest": 10,
        "loanWeeks": 50
    }'

    ```

    Response:
    ```json
    {
        "billingId": 1,
        "outstanding": 5500000,
        "customerId": 1,
        "loanId": 1001,
        "loanAmount": 5000000,
        "loanInterest": 10,
        "loanWeeks": 50
    }
    ```

- Get Billing
    Request:
    ```curl
    curl -X GET http://localhost:8080/api/v1/billings/1
    ```

    Response:
    ```json
    {
        "id": 1,
        "customerId": 1,
        "loanId": 1001,
        "loanAmount": 5000000,
        "loanWeeks": 50,
        "loanInterest": 10,
        "outstanding": 5500000,
        "Payments": [
            {
                "id": 1,
                "billingId": 1,
                "amount": 110000,
                "week": 1,
                "paid": false,
                "CreatedAt": "2025-08-06T09:27:29.140612Z",
                "UpdatedAt": "2025-08-06T09:27:29.140612Z",
                "DeletedAt": null
            },
            {
                "id": 2,
                "billingId": 1,
                "amount": 110000,
                "week": 2,
                "paid": false,
                "CreatedAt": "2025-08-06T09:27:29.140612Z",
                "UpdatedAt": "2025-08-06T09:27:29.140612Z",
                "DeletedAt": null
            },
            {
                "id": 3,
                "billingId": 1,
                "amount": 110000,
                "week": 3,
                "paid": false,
                "CreatedAt": "2025-08-06T09:27:29.140612Z",
                "UpdatedAt": "2025-08-06T09:27:29.140612Z",
                "DeletedAt": null
            },
            ...
            {
                "id": 50,
                "billingId": 1,
                "amount": 110000,
                "week": 50,
                "paid": false,
                "CreatedAt": "2025-08-06T09:27:29.140612Z",
                "UpdatedAt": "2025-08-06T09:27:29.140612Z",
                "DeletedAt": null
            }
        ],
        "CreatedAt": "2025-08-06T09:27:29.127608Z",
        "UpdatedAt": "2025-08-06T09:27:29.127608Z",
        "DeletedAt": null
    }
    ```

- Get Outstanding
    Request:
    ```curl
    curl -X GET http://localhost:8080/api/v1/billings/1/outstanding
    ```

    Response:
    ```json
    {
        "billingId": 1,
        "customerId": 1,
        "loanId": 1001,
        "outstanding": 5500000
    }
    ```

- Isdelinquent
    Request:
    ```curl
    curl -X GET http://localhost:8080/api/v1/billings/1/delinquent
    ```

    Response:
    ```json
    {
        "billingId": 1,
        "customerId": 1,
        "loanId": 1001,
        "isDelinquent": true
    }
    ```
- Make Payment
Request:
    ```curl
    curl -X POST http://localhost:8080/api/v1/billings/1/payments \
    -H "Content-Type: application/json" \
    -d '{
        "week": 1,
        "amount": 110000
    }'
    ```

    Response:
    ```json
    {
        "customerId": 1,
        "loanId": 1001,
        "outstanding": 5390000,
        "payment": {
            "id": 1,
            "billingId": 1,
            "amount": 110000,
            "week": 1,
            "paid": true,
            "CreatedAt": "2025-08-06T09:27:29.140612Z",
            "UpdatedAt": "2025-08-06T10:06:04.982112354Z",
            "DeletedAt": null
        }
    }
    ```