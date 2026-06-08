import { createRouter, createWebHistory } from "vue-router";

import LoginView from "../views/LoginView.vue";
import DashboardView from "../views/DashboardView.vue";
import AssetsView from "../views/AssetsView.vue";
import IncidentsView from "../views/IncidentsView.vue";
import UsersView from "../views/UsersView.vue";
import NotificationsView from "../views/NotificationsView.vue";
import SitesView from "../views/SitesView.vue";
import ServicesView from "../views/ServicesView.vue";
import ApplicationsView from "../views/ApplicationsView.vue";
import DatabasesView from "../views/DatabasesView.vue";
import RelationsView from "../views/RelationsView.vue";
import ChangesView from "../views/ChangesView.vue";
import CmdbGraphView from "../views/CmdbGraphView.vue";
import ChangeImpactView from "../views/ChangeImpactView.vue";

const routes = [
    { path: "/", component: LoginView },
    { path: "/dashboard", component: DashboardView },
    { path: "/assets", component: AssetsView },
    { path: "/incidents", component: IncidentsView },
    { path: "/users", component: UsersView },
    { path: "/sites", component: SitesView },
    { path: "/services", component: ServicesView },
    { path: "/notifications", component: NotificationsView },
    { path: "/cmdb/applications", component: ApplicationsView },
    { path: "/cmdb/databases", component: DatabasesView },
    { path: "/cmdb/relations", component: RelationsView },
    { path: "/cmdb/graph", component: CmdbGraphView },
    { path: "/changes", component: ChangesView },
    { path: "/changes/:id/impact", component: ChangeImpactView},
    { path: "/changes/:id/impact", name: "change-impact", component: ChangeImpactView},

];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

export default router;