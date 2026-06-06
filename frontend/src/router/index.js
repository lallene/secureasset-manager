import { createRouter, createWebHistory } from "vue-router";

import LoginView from "../views/LoginView.vue";
import DashboardView from "../views/DashboardView.vue";
import AssetsView from "../views/AssetsView.vue";
import IncidentsView from "../views/IncidentsView.vue";
import UsersView from "../views/UsersView.vue";
import NotificationsView from "../views/NotificationsView.vue";
import SitesView from "../views/SitesView.vue";
import ServicesView from "../views/ServicesView.vue";

const routes = [
    { path: "/", component: LoginView },
    { path: "/dashboard", component: DashboardView },
    { path: "/assets", component: AssetsView },
    { path: "/incidents", component: IncidentsView },
    { path: "/users", component: UsersView },
    { path: "/sites", component: SitesView },
    { path: "/services", component: ServicesView },
    { path: "/notifications", component: NotificationsView },
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;