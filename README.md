# Webring

Software to create a webring

## Usage

To create your own webring, you have to install the application by building it by yourself or via Docker.

With this, you can print an example config to the std out with the flag `-generate-config`.
Its format is the TOML.

Then, you can start your server with the flag `-config path/to/config.toml`.

### HTML files
You must create two HTML files for the pages join and legal information.
Their path is set in the config file.

You can use all common tags (like `p`, `a`, `h2`, `img`...).
The CSS class `text` is used when you want to display text.
The class `list` is used when you want to display a list into columns (like in the home).

### Public files
All static files are located inside the `public` directory.
These are available under the URI `/static/path_in_public`.

You can modify the favicon by putting an image called `logo.png` in your public directory.
You can change this name in the config.

## Technologies

It is an SSR compiled webserver.

- Go 1.24
- [anhgelus/Golatt](https://github.com/anhgelus/golatt)
- Bun
- SCSS
