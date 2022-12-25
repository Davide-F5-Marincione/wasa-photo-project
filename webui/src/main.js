import {createApp, reactive} from 'vue'
import { Buffer } from 'buffer';
globalThis.Buffer = Buffer;
import JSONbigint from 'json-bigint' //The reason for this is described below!
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import TopBar from './components/TopBar.vue'


import './assets/topbar.css'
import './assets/posts.css'
import './assets/searchbar.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;

axios.interceptors.request.use(
    config => {
        // For whatever reason, (I don't know why and, frankly, I don't want to)
        // axios is unable to keep storing the authorization parameter once I set it.
        // Therefore the solution (except manually setting the parameter at every call) is to
        // just modify the problem here by resetting it manually at this point: it's not good-looking
        // but it works.
        // See more here https://stackoverflow.com/questions/55137517/bearer-token-undefined-with-global-axios-config
        const token = localStorage.getItem('token');
        const auth = token ? `Bearer ${token}` : '';
        config.headers.common['Authorization'] = auth;
        return config;
    }
);


// I'm using 64bit integers for the authorization
// and during the login this authorization is defined to be returned
// as an actual integer.
// Unfortunately though axios uses JSON.parse (or whatever) and because
// of that it doesn't handle big integers well (puts a bunch of zeroes once
// it doesn't want to parse the integer anymore).
// Since I don't want to fix the issue at the source by modifying my API and backend
// (also because this is bullshit, why the hell don't you parse up to 64bit at least?!)
// I'm manually replacing JSON.parse to JSONBigInt.parse, which is able to parse
// big integers...
// See more here https://stackoverflow.com/questions/43787712/axios-how-to-deal-with-big-integers
// (BTW to debug this was hellish, how could I know something like this was a problem?)
axios.interceptors.response.use(response => {
    if (response.headers["content-type"] == "application/json"){
        response.data = JSONbigint.parse(response.data)
    }
    return response
})

app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("TopBar", TopBar);
app.component("Buffer", Buffer);
app.use(router)
app.mount('#app')
