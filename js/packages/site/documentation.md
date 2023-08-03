# Getting Started

## Prerequisites

Makes sure you have golang installed

<https://go.dev>

## Installation

1. Install server

   ```sh
   go install github.com/go-uvid/uvid@latest
   ```

2. Run server

   ```sh
   go/bin/uvid
   ```

3. Open dashboard with `YOUR_SERVER_IP/:3000`, or <http://localhost:3000> if you installed locally
4. Login to dashboard, default user name is `root`, password is `uvid`
5. That's all. `./uvid.db` is sqlite database stores all your data

## Usage

To start track your website, you must install SDK in your front-end project.

1. Install SDK

```sh
npm install uvid-js
```

2. Initialize SDK, By default, SDK will sends page view and performance data on page load

```js
import { init } from "uvid-js";

const sdk = init({
  host: "YOUR_SERVER_IP/:3000",
  sessionMeta: {
    // optionally, track additional meta data
    userId: "123",
  },
  // optionally, provide your website's build version
  appVersion: "1.0.0",
});

// Track an js error
sdk.error(new Error("This is an js error!"));
// Track an custom event action and value
sdk.event("register", "some-user@email.com");
// Track an custom event by HTML attributes
`<button data-uvid-action="register" data-uvid-value="some-user@email.com">Register</button>`;
//  When user click the button, uvid-js will track it and call uvid.event('register', 'some-user@email.com')

// Track an http request
sdk.http({
  resource: "http://some-api.com",
  status: 500,
  method: "GET",
  headers: "",
});
```

3. Publish above change, when user visit your website, you will see data from dashboard
