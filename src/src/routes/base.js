
const Error404 = {
  path: '*',
  name: 'error_404',
  component: () => import('@v/error-page/error-404.vue'),
  meta: {
    title: '404',
  },
};
export {
  Error404,
};
