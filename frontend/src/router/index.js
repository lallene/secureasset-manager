import { createRouter, createWebHistory } from "vue-router";

import LoginView from "../views/LoginView.vue";
import DashboardView from "../views/DashboardView.vue";
import AssetsView from "../views/AssetsView.vue";
import IncidentsView from "../views/IncidentsView.vue";
import UsersView from "../views/UsersView.vue";

const routes = [
    { path: "/", component: LoginView },
    { path: "/dashboard", component: DashboardView },
    { path: "/assets", component: AssetsView },
    { path: "/incidents", component: IncidentsView },
    { path: "/users", component: UsersView },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;