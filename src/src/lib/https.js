import HttpRequest from './request';
import config from '@/config';
const apiUrl = location.origin.indexOf('mct.api') > -1 ? `${location.origin}/index.php/web/` : `${location.origin}/merchant/index.php/web/`;
const baseUrl = process.env.NODE_ENV === 'development' ? config.baseUrl.dev : apiUrl;
const https = new HttpRequest(baseUrl);
export default https;
