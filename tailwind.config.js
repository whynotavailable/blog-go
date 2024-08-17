/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./templates/*.templ"
  ],
  safelist: [
    'size-2/5'
  ],
  theme: {
    extend: {
      fontFamily: {
        body: '"Open Sans", sans-serif',
        mono: '"Roboto Mono", monospace'
      },
      colors: {
        black: '#161616',
        white: '#f2f4f8',
        on: '#f16da6',
        code: '#c8a5ff',
        headers: '#8cb6ff',
        muted: '#52bdff'
      }
    },
  },
  plugins: [],
}

