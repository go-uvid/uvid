<a name="readme-top"></a>

[![Stargazers][stars-shield]][stars-url]
[![MIT License][license-shield]][license-url]

<!-- PROJECT LOGO -->
<br />
<div align="center">
  <a href="[repo-url]">
    <img src="media/logo.svg" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">UVID</h3>

  <p align="center">
    <br />
    Observable Platform for Frontend Websites
    <br />
    <br />
<!-- <a href="https://github.com/rick-you/uvid">View Demo</a>
    · -->
    <a href="[issues-url]">Report Bug</a>
·
    <a href="[issues-url]">Request Feature</a>
  </p>
</div>

<!-- ABOUT THE PROJECT -->

## About The Project

UVID can help you:

- Tracks real user interactions
- Monitors site errors
- Captures site performance

<!-- [![Product Name Screen Shot][product-screenshot]](https://example.com) -->

### Built With

- ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
- ![SQLite](https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white)
- ![React](https://img.shields.io/badge/react-%2320232a.svg?style=for-the-badge&logo=react&logoColor=%2361DAFB)

## Getting Started

### Prerequisites

Makes sure you have golang installed

<https://go.dev>

### Installation

1. Install server

   ```sh
   go install github.com/rick-you/uvid@latest
   ```

2. Run server

   ```sh
   uvid
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
import { init } from 'uvid-js';

const sdk = init({
  host: 'YOUR_SERVER_IP/:3000',
  sessionMeta: {
    // optionally, track additional meta data
    userId: '123',
  },
  // optionally, provide your website's build version
  appVersion: '1.0.0',
});

// Track an js error
sdk.error(new Error('This is an js error!'));
// Track an custom event action and value
sdk.event('register', 'some-user@email.com');
// Track an custom event by HTML attributes
`<button data-uvid-action="register" data-uvid-value="some-user@email.com">Register</button>`;
//  When user click the button, uvid-js will track it and call uvid.event('register', 'some-user@email.com')

// Track an http request
sdk.http({
  resource: 'http://some-api.com',
  status: 500,
  method: 'GET',
  headers: '',
});
```

3. Publish above change, when user visit your website, you will see data from dashboard

## Contributing

TODO

## License

See `LICENSE.txt` for more information.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->

[stars-shield]: https://img.shields.io/github/stars/rick-you/uvid.svg?style=for-the-badge
[stars-url]: https://github.com/rick-you/uvid/stargazers
[license-shield]: https://img.shields.io/github/license/rick-you/uvid.svg?style=for-the-badge
[license-url]: https://github.com/rick-you/uvid/blob/master/LICENSE.txt
[issues-url]: https://github.com/rick-you/uvid/issues
