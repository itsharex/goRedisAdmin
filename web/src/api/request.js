//对于axios进行二次封装
import axios from "axios";
import {Message} from 'element-ui'
import {GetToken} from "@/utils/token";

//底下的代码也是创建axios实例
let requests = axios.create({
    //基础路径
    baseURL: "/api/v1/",
    timeout: 5000,
});

//请求拦截器----在项目中发请求（请求没有发出去）可以做一些事情
requests.interceptors.request.use((config) => {
    config.headers['Authorization'] = 'Bearer ' + GetToken()
    return config
});

//响应拦截器----当服务器手动请求之后，做出响应（相应成功）会执行的
requests.interceptors.response.use(
    response => {
        // 判断如果是下载请求直接返回
        let rt = response.config.data
        console.log(response)
        if (rt != undefined) {
            rt = JSON.parse(rt)
            if (rt.responseType == "blob") {
                return response
            }
        }
        let res = response.data
        if (res.code !== 0 && res.code !== 200) {
            if (res.code == 6) {
                Message({
                    message: '登录无效，请重新登录',
                    type: 'error',
                    duration: 5 * 1000
                })
                window.location.href = '/#/login'; // 或者使用路由导航实现跳转
            } else {
                Message({
                    message: res.message || 'Error',
                    type: 'error',
                    duration: 5 * 1000
                })
            }
            return Promise.reject(new Error(res.message || 'Error'))
        } else {
            return res
        }
    },
    error => {
        console.log('err' + error) // for debug
        Message({
            message: error.message,
            type: 'error',
            duration: 5 * 1000
        })
        return Promise.reject(error)
    }
);
//最终需要对外暴露（不对外暴露外面模块没办法使用）
//这里的代码是暴露一个axios实例
export default requests;
