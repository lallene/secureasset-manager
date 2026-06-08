<script setup>
import { onMounted, ref } from "vue";
import api from "../api/axios";
import DashboardLayout from "../layouts/DashboardLayout.vue";

const services = ref([]);
const showModal = ref(false);
const isEditMode = ref(false);
const editingServiceId = ref(null);
const error = ref("");
const isSubmitting = ref(false);

const form = ref({
  name: "",
  description: "",
});

const getToken = () => localStorage.getItem("token");

const fetchServices = async () => {
  try {
    const response = await api.get("/services", {
      headers: { Authorization: `Bearer ${getToken()}` },
    });
    services.value = response.data;
  } catch (err) {
    console.error(err);
    error.value = "Impossible de charger les services";
  }
};

const resetForm = () => {
  form.value = {
    name: "",
    description: "",
  };
  isEditMode.value = false;
  editingServiceId.value = null;
};

const openCreateModal = () => {
  resetForm();
  showModal.value = true;
};

const openEditModal = (service) => {
  isEditMode.value = true;
  editingServiceId.value = service.ID;

  form.value = {
    name: service.name,
    description: service.description,
  };

  showModal.value = true;
};

const saveService = async () => {
  try {
    isSubmitting.value = true;
    error.value = "";

    if (isEditMode.value) {
      await api.put(`/services/${editingServiceId.value}`, form.value, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    } else {
      await api.post("/services", form.value, {
        headers: { Authorization: `Bearer ${getToken()}` },
      });
    }

    showModal.value = false;
    resetForm();
    await fetchServices();
  } catch (err) {
    console.error(err);
    error.value = "Impossible d'enregistrer le service";
  } finally {
    isSubmitting.value = false;
  }
};

const deleteService = async (service) => {
  if (!confirm(`Supprimer le service "${service.name}" ?`)) return;

  try {
    await api.delete(`/services/${service.ID}`, {
      headers: { Authorization: `Bearer ${getToken()}` },
    });

    await fetchServices();
  } catch (err) {
    console.error(err);
    error.value = "Impossible de supprimer le service";
  }
};

onMounted(fetchServices);
</script>

<template>
  <DashboardLayout>
    <!-- Topbar Header -->
    <div class="topbar flex items-center justify-between">
      <div>
        <h1>Services</h1>
        <p>Gérez les services responsables du traitement des tickets.</p>
      </div>

      <button @click="openCreateModal" class="btn-primary flex items-center gap-2">
        <span>+</span> Ajouter un service
      </button>
    </div>

    <!-- Alert Bar -->
    <div v-if="error" class="modal-error error-bar mb-6">
      {{ error }}
    </div>

    <!-- Main Table Container -->
    <div class="table-card overflow-hidden">
      <table class="w-full text-left border-collapse">
        <thead>
        <tr>
          <th class="p-4 text-xs font-semibold uppercase tracking-wider">Nom</th>
          <th class="p-4 text-xs font-semibold uppercase tracking-wider">Description</th>
          <th class="p-4 text-xs font-semibold uppercase tracking-wider text-right">Actions</th>
        </tr>
        </thead>

        <tbody>
        <tr v-for="service in services" :key="service.ID">
          <td class="p-4 font-medium" style="color: var(--tx-primary);">
            {{ service.name }}
          </td>
          <td class="p-4">
            {{ service.description || "-" }}
          </td>
          <td class="p-4 text-right">
            <div class="flex justify-end gap-2">
              <button @click="openEditModal(service)" class="action-btn action-edit">
                Modifier
              </button>
              <button @click="deleteService(service)" class="action-btn action-delete">
                Supprimer
              </button>
            </div>
          </td>
        </tr>

        <!-- Empty State intégré aux standards de ton CSS -->
        <tr v-if="services.length === 0">
          <td colspan="3" class="p-12 text-center">
            <div class="flex flex-col items-center justify-center gap-2">
              <div class="empty-title font-medium text-base">Aucun service trouvé</div>
              <div class="empty-sub text-sm">Commencez par ajouter un nouveau service à la liste.</div>
            </div>
          </td>
        </tr>
        </tbody>
      </table>
    </div>
  </DashboardLayout>

  <!-- Modal Component -->
  <div v-if="showModal" class="modal-overlay fixed inset-0 flex items-center justify-center p-4 z-50">
    <div class="modal rounded-2xl w-full max-w-lg overflow-hidden shadow-2xl">

      <div class="modal-header p-5 flex justify-between items-center border-b">
        <div>
          <span class="modal-eyebrow text-xs font-bold uppercase tracking-wider block mb-1">Configuration</span>
          <h2 class="modal-title text-lg font-semibold m-0">
            {{ isEditMode ? "Modifier le service" : "Ajouter un service" }}
          </h2>
        </div>
        <button @click="showModal = false" class="modal-close p-2 rounded-lg transition-colors">
          ✕
        </button>
      </div>

      <form @submit.prevent="saveService" class="p-6 space-y-5">
        <div class="flex flex-col gap-1.5">
          <label class="field-label text-xs font-medium">Nom du service</label>
          <input
              v-model="form.name"
              required
              placeholder="Ex: Support Technique, Facturation..."
              class="field-input w-full"
          />
        </div>

        <div class="flex flex-col gap-1.5">
          <label class="field-label text-xs font-medium">Description</label>
          <textarea
              v-model="form.description"
              placeholder="Décrivez brièvement le périmètre de ce service..."
              rows="4"
              class="field-input w-full resize-none"
          ></textarea>
        </div>

        <div class="modal-footer -mx-6 -mb-6 p-4 mt-6 flex justify-end gap-3 border-t">
          <button type="button" @click="showModal = false" class="btn-ghost">
            Annuler
          </button>
          <button type="submit" :disabled="isSubmitting" class="btn-primary min-w-[100px]">
            <span v-if="isSubmitting" class="inline-block animate-spin mr-2">✦</span>
            {{ isEditMode ? "Enregistrer" : "Créer" }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
/* Alignement et ajustements spécifiques au layout des boutons internes */
.btn-primary {
  padding: 10px 18px;
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--r-md);
  cursor: pointer;
  transition: background var(--t-fast);
}

.btn-ghost {
  padding: 10px 18px;
  font-size: 14px;
  font-weight: 500;
  border-radius: var(--r-md);
  cursor: pointer;
  transition: background var(--t-fast);
}

/* Boutons d'action en bout de ligne */
.action-btn {
  padding: 6px 12px;
  font-size: 13px;
  font-weight: 500;
  border-radius: var(--r-sm);
  border: 1px solid transparent;
  cursor: pointer;
  transition: all var(--t-fast);
}

/* Styles pour la barre d'erreur */
.error-bar {
  padding: 14px 16px;
  border-radius: var(--r-md);
  border: 1px solid;
  font-size: 14px;
}

/* Modal header/footer tweaks */
.modal-header,
.modal-footer {
  border-style: solid;
}
</style>