## Authentication API

This API is responsible for generating and validating JWT tokens. Some of the functions of this API include:

* When a user first logs into the system, a request is made to this API which then verifies the credentials with the Account API. Once this validation is passed, a jwt is generated which contains information about the user.

* jwt is used throughout the user journey to verify its originality. Examples include user session, creating new records, accessing site content, modifying records.

* each jwt is signed by a private key which is stored on the server