module.exports = {
  root: true,
  env: {
    node: true,
  },
  'extends': [
    'plugin:vue/essential',
    '@vue/standard',
  ],
  globals: {
    'moment': true,
    'antd': true,
  },
  rules: {
    'comma-dangle': [0, 'always',],
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',
    'semi': ["error", "always",], // 强行加分号
    'space-before-function-paren': ["error", "never",],
    'quotes': 'off',
  },
  parserOptions: {
    parser: 'babel-eslint',
  },
};
