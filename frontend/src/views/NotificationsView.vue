<script setup>
import { onMounted, ref, computed } from "vue";
import { useRouter } from "vue-router";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const notifications = ref([]);
const error = ref("");

const getToken = () => localStorage.getItem("token");

const unreadCount = computed(() => {
  return notifications.value.filter((n) => !n.is_read).length;
});

const router = useRouter();

const goToIncident = async (notification) => {
  if (!notification.incident_id) return;
  if (!notification.is_read) {
    await markAsRead(notification);
  }
  router.push(`/incidents?incident=${notification.incident_id}`);
};

const fetchNotifications = async () => {
  try {
    const response = await api.get("/notifications", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    notifications.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les notifications";
  }
};

const markAsRead = async (notification) => {
  try {
    await api.put(`/notifications/${notification.ID}/read`, {}, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    notification.is_read = true;
  } catch (err) {
    console.error(err);
  }
};

const markAllAsRead = async () => {
  const unreadList = notifications.value.filter(n => !n.is_read);
  for (const n of unreadList) {
    await markAsRead(n);
  }
};

onMounted(fetchNotifications);
</script>

<template>
  <DashboardLayout>
    <div class="topbar flex items-center justify-between">
      <div>
        <h1>Notifications</h1>
        <p>Restez informé des mises à jour sur vos incidents et vos infrastructures.</p>
      </div>

      <div class="flex items-center gap-3">
        <button
            v-if="unreadCount > 0"
            @click="markAllAsRead"
            class="btn-ghost text-xs"
            style="padding: 8px 14px;"
        >
          Tout marquer comme lu
        </button>
        <span
            class="sc-bdg px-3 py-1 rounded-full text-xs font-semibold border"
            :class="unreadCount > 0 ? 'bdg-brand' : 'bdg-sl'"
        >
          {{ unreadCount }} non lues
        </span>
      </div>
    </div>

    <div v-if="error" class="modal-error error-bar mb-6 flex items-center gap-2">
      <span class="w-1.5 h-1.5 rounded-full bg-rose-500 pulse"></span>
      {{ error }}
    </div>

    <div class="space-y-4 max-w-4xl">
      <div
          v-for="notification in notifications"
          :key="notification.ID"
          class="nt-card transition-all duration-200 p-5 rounded-2xl border relative"
          :class="!notification.is_read ? 'nt-unread' : 'nt-read'"
      >
        <div class="flex items-start gap-4">
          <div class="mt-1.5 flex-shrink-0">
            <span
                class="block w-2.5 h-2.5 rounded-full transition-all duration-300"
                :class="!notification.is_read ? 'bg-brand shadow-glow' : 'bg-muted'"
            ></span>
          </div>

          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-4">
              <div>
                <h2
                    class="text-base m-0 transition-colors"
                    :style="{ color: !notification.is_read ? 'var(--tx-primary)' : 'var(--tx-muted)' }"
                    :class="!notification.is_read ? 'font-semibold' : 'font-medium'"
                >
                  {{ notification.title }}
                </h2>
                <p class="text-sm mt-1.5 leading-relaxed" style="color: var(--tx-secondary);">
                  {{ notification.message }}
                </p>
              </div>

              <span class="text-xs whitespace-nowrap pt-0.5" style="color: var(--tx-muted);">
                {{ new Date(notification.CreatedAt).toLocaleString(undefined, { hour: '2-digit', minute:'2-digit', day: 'numeric', month: 'short' }) }}
              </span>
            </div>

            <div class="flex items-center gap-3 mt-4 pt-4 border-t" style="border-color: var(--bd-subtle);">
              <button
                  v-if="notification.incident_id"
                  @click="goToIncident(notification)"
                  class="action-btn action-edit flex items-center gap-1.5 text-xs"
              >
                <span>Ouvrir l'incident</span>
                <span class="text-xs">➔</span>
              </button>

              <button
                  v-if="!notification.is_read"
                  @click="markAsRead(notification)"
                  class="btn-ghost text-xs"
                  style="padding: 6px 12px; border-radius: var(--r-sm);"
              >
                Marquer comme lu
              </button>
            </div>
          </div>
        </div>
      </div>

      <div v-if="notifications.length === 0" class="table-card p-12 text-center border-dashed max-w-4xl">
        <div class="w-12 h-12 rounded-full flex items-center justify-center mx-auto mb-3" style="background: var(--bg-base); color: var(--tx-muted);">
          🔔
        </div>
        <div class="empty-title font-medium text-base">Aucune notification</div>
        <div class="empty-sub text-sm mt-1">Vous êtes complètement à jour ! Tout fonctionne correctement.</div>
      </div>
    </div>
  </DashboardLayout>
</template>

<style scoped>
/* Boutons & Composants Génériques du Thème */
.btn-primary, .btn-ghost {
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--r-md);
  cursor: pointer;
  transition: all var(--t-fast);
}

.action-btn {
  padding: 6px 14px;
  font-weight: 500;
  border-radius: var(--r-sm);
  cursor: pointer;
  transition: all var(--t-fast);
}

.error-bar {
  padding: 14px 16px;
  border-radius: var(--r-md);
  border: 1px solid;
  font-size: 14px;
}

/* Styles spécifiques aux conteneurs de notification */
.nt-card {
  background: var(--bg-card);
  border-color: var(--bd-subtle);
}

.nt-unread {
  border-color: var(--brand-subtle, var(--brand));
  background: linear-gradient(90deg, rgba(127, 119, 221, 0.04) 0%, rgba(0, 0, 0, 0) 100%), var(--bg-card);
}

.nt-read {
  opacity: 0.8;
}

.bg-brand {
  background-color: var(--brand);
}

.bg-muted {
  background-color: var(--bd-subtle);
}

.shadow-glow {
  box-shadow: 0 0 8px var(--brand);
}
</style>