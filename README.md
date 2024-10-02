# `ktz` - Timezone CLI Application

`ktz` is a simple and efficient timezone CLI application that helps you retrieve the current time based on a given city, country, or timezone. It's perfect for quickly finding the time anywhere in the world through your terminal.

---

## Table of Contents

1. [Installation](#installation)
2. [Usage / Examples](#usage-examples)
3. [Future Plans](#future-plans)

---

## 1. Installation

To get started with `ktz`, follow either <strong>one</strong> of the steps below:

1. Install the application by running:

   ```zsh
   $ go install github.com/kritibb/ktz@v0.1.0
   ```

2. Or, run the following commands:

   ```zsh
   $ git clone git@github.com:kritibb/ktz.git
   $ cd ktz
   $ sh setup_ktz.sh
   ```

## 2. Usage / Examples

#### Find Timezone by City

- Use full name of a city:

  ```bash
  $ ktz lookup sydney
  ```

- Use prefix search:

  ```bash
  $ ktz lookup syd
  ```

    <details>
    <summary><strong>Demo</strong></summary>

  ![City prefix-search demo](https://github.com/kritibb/documentation-files/blob/main/ktz-assets/city-trie.gif?raw=true)

    </details>

#### Find Timezone by Country

- Use the 3-letter country code:

  ```bash
  $ ktz lookup -c NPL
  ```

    <details>
    <summary><strong>Demo</strong></summary>

  ![Country code demo](https://github.com/kritibb/documentation-files/blob/main/ktz-assets/countrycode3.gif?raw=true)

    </details>

- Use the 3-letter country code:

  ```bash
  $ ktz lookup -c NP
  ```

- Use the prefix search:

  ```bash
  $ ktz lookup -c aust
  ```

    <details>
    <summary><strong>Demo</strong></summary>

  ![Country prefix-search demo](https://github.com/kritibb/documentation-files/blob/main/ktz-assets/country.gif?raw=true)

    </details>

#### Find Timezone by tz name

- Use timezone abbreviation:

  ```bash
  $ ktz lookup -z pst
  ```

    <details>
    <summary><strong>Demo</strong></summary>

  ![Zone abbreviation demo](https://github.com/kritibb/documentation-files/blob/main/ktz-assets/zone-abb.gif?raw=true)

    </details>

- Use the full timezone name:

  ```bash
  $ ktz lookup -z Asia/Kathmandu
  ```

    <details>
    <summary><strong>Demo</strong></summary>

  ![Zone name demo](https://github.com/kritibb/documentation-files/blob/main/ktz-assets/zone-full.gif?raw=true)

    </details>

## 3. Future Plans

Here are some features that we plan to add in the future:

- Favorite Timezones: Add specific timezones to your favorite list for quick access.
- Remove Timezones: Remove timezones from your favorite list.
- Display Saved Timezones: Easily view all your saved timezones in one place.

Stay tuned for more updates!

Feel free to contribute, suggest improvements, or raise issues to help us make `ktz` even better! 😊
