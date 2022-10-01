const axios = require("axios").default;
import Router from "next/router";

const instance = axios.create({
  baseURL: process.env.NEXT_PUBLIC_API_SERVER_BASE_URI,
  withCredentials: true,
  timeout: 10000,
});
instance.interceptors.response.use(
  function (response) {
    return response;
  },
  function (error) {
    if (error.response.status == 401) {
      Router.push("/");
    }
    return Promise.reject(error);
  }
);

export function getUserSettings() {
  return instance
    .get("/user/settings")
    .then((response) => response.data)
    .catch(function (error) {});
}
