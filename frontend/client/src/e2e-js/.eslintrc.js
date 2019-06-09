module.exports = {
  'env': {
    'browser': true,
    'commonjs': true,
    'node': true,
    'es6': true
  },
  'parserOptions': {
    'ecmaVersion': 8
  },
  'extends': 'eslint:recommended',
  'rules': {
    'linebreak-style': [
      'error',
      'unix'
    ],
    'quotes': [
      'error',
      'single'
    ],
    'semi': [
      'error',
      'always'
    ]
  }
};
