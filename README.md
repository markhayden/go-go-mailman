# Go Go Mailman

## Overview
A simple transactional email rendering layer. Sits between the trigger and the actual SMTP. This service takes a POST request, renders the email template for the request and then passes the full request on to the SMTP (mailgun at the moment). Does some simple send logging but nothing fancy at the moment.

## Templates
Templates by default are managed in the `tmpl` directory. This directory can be reconfigured with ease using the `config.yaml` file.

## Details
**Type:** POST  
**Endpoint:** /deploy/email  

**Params:**  

| Variable | Value | Required | Notes |
|---|---|---|---|
|Token|String|True|
|Template|String|True|Represents the file name without type.|
|Params|Object|False|{"food":"taco", "topping":"lettuce"}|
|Recipient|Array of Objects|True|[{"name":"Mark Hayden", "email":"my@email.com"}]|
|Subject|String|True|


## Example
```
echo '{
    "token": "14l2kh4g2erhg2345yh2",
    "template": "test",
    "params": {
        "test": "one",
        "test2": "2"
    },
    "recipient": [
        {
            "name": "Mark",
            "email": "imanemail@email.com"
        }
    ],
    "subject": "test"
}' | http POST localhost:8080/deploy/email
```