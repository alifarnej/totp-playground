# Time-Based One-Time Password (TOTP)

## What is TOTP?

TOTP (Time-Based One-Time Password) is a temporary, one-time password algorithm that generates a unique password which is valid only for a short period of time. It is widely used for two-factor authentication (2FA) to enhance security by requiring two forms of verification.

## How TOTP Generates Codes

TOTP generates a unique code based on two key components:
1. **Shared Secret**: A secret key shared between the server and the client (e.g., a user’s mobile app).
2. **Current Time**: The current timestamp, typically in 30-second intervals.

The steps to generate a TOTP code are as follows:

1. **Generate a HMAC-SHA-1 Hash**: Combine the shared secret with the current timestamp (number of 30-second intervals since the Unix epoch).
2. **Truncate the Hash**: Extract a dynamic binary code from the hash.
3. **Convert to Decimal**: Convert the binary code to a 6-digit decimal number (the TOTP code).

## TOTP Code Validation

To validate a TOTP code, the server performs the following steps:
1. **Retrieve Shared Secret**: Retrieve the shared secret associated with the user.
2. **Calculate Expected TOTP Code**: Using the same algorithm and the current timestamp, calculate the expected TOTP code.
3. **Compare Codes**: Compare the user-provided TOTP code with the expected code. If they match, the authentication is successful.

## Real Example

### Step 1: Shared Secret

Let's assume the shared secret between the server and client is `JBSWY3DPEHPK3PXP`.
For this example, let's assume the Base32 decoding results in the following byte array:

```shell
[0x4a, 0x42, 0x53, 0x57, 0x59, 0x33, 0x44, 0x50, 0x45, 0x48, 0x50, 0x4b, 0x33, 0x50, 0x58, 0x50]
```

### Step 2: Current Time

Assume the current timestamp is `2024-08-05 10:00:00 UTC`. The Unix time for this is `1722825600` seconds since the epoch. 
We then divide this by 30 to get the time step:

```shell
T = 1722825600 / 30 = 57427520
```
The time step `57427520` needs to be represented as an 8-byte integer in big-endian format. 
Suppose formatting is done as follows:
```shell
Time Step: 57427520
Time Step Bytes: [0x00, 0x00, 0x00, 0x00, 0x03, 0x6e, 0x80, 0x00]
```

### Step 3: Generate HMAC-SHA-1 Hash

Use HMAC with the SHA-1 hashing algorithm. The input for HMAC-SHA-1 is:
- Key: Shared Secret Bytes
- Message: Time Step Bytes

For this example, let's assume the resulting hash is:

```shell
HMAC-SHA-1 Hash: [0x4a, 0x6c, 0x1f, 0x48, 0x9f, 0x9f, 0xe3, 0x79, 0x1b, 0x62, 0x6a, 0xd2, 0x0f, 0x6b, 0x5a, 0x64, 0x1f, 0xd7, 0x64, 0x14]
```

### Step 4: Dynamic Truncation

Extract a 4-byte dynamic binary code from the hash:

```shell
Offset = Hash[19] & 0xf = 0x4
4-byte Code = Hash[4:8] = [0x9f, 0x9f, 0xe3, 0x79]
```

### Step 5: Convert to TOTP Code

Convert the binary code to a decimal TOTP code, lets assume:

```shell
Binary Code = 0x2663ef70 = 644808432
TOTP Code = 644808432 % 10^6 = 808432
```


### Step 6: Validation

When the user enters the TOTP code `808432`, the server will generate the expected code using the same shared secret and current timestamp. If the code matches, the authentication is successful.


