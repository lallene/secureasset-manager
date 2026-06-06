<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";
import api from "../api/axios";

const router = useRouter();

const email = ref("");
const password = ref("");
const error = ref("");
const isLoading = ref(false); // Ajout d'un état de chargement pour l'élégance et l'UX

const login = async () => {
  if (!email.value || !password.value) {
    error.value = "Veuillez remplir tous les champs";
    return;
  }

  try {
    error.value = "";
    isLoading.value = true;

    const response = await api.post("/login", {
      email: email.value,
      password: password.value,
    });

    localStorage.setItem("token", response.data.token);
    localStorage.setItem("role", response.data.role);
    localStorage.setItem("name", response.data.name);

    router.push("/dashboard");
  } catch (err) {
    error.value = "Email ou mot de passe incorrect";
  } finally {
    isLoading.value = false;
  }
};
</script>

<template>
  <!-- Fond adouci avec un très léger dégradé radial pour donner de la profondeur -->
  <div class="min-h-screen flex flex-col items-center justify-center bg-slate-950 bg-[radial-gradient(ellipse_at_top,_var(--tw-gradient-stops))] from-slate-900 via-slate-950 to-slate-950 p-4 font-sans text-slate-900 antialiased">

    <div class="w-full max-w-md">
      <!-- Logo / Branding au-dessus de la carte -->
      <div class="text-center mb-8">
        <div class="inline-flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-tr from-blue-600 to-indigo-500 shadow-lg shadow-blue-500/20 mb-4">
          <!-- Icône minimaliste de bouclier/clé (SVG Fin) -->
          <svg class="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke-width="2" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z" />
          </svg>
        </div>
        <h1 class="text-2xl font-bold tracking-tight text-white">SecureAsset Manager</h1>
        <p class="text-sm text-slate-400 mt-1.5">Connectez-vous pour gérer vos infrastructures</p>
      </div>

      <!-- Carte du formulaire -->
      <div class="bg-white rounded-2xl shadow-xl shadow-black/40 border border-slate-100 p-8">
        <!-- Formulaire natif (permet la soumission via la touche Entrée) -->
        <form @submit.prevent="login" class="space-y-5">

          <!-- Champ Email -->
          <div>
            <label for="email" class="block text-xs font-semibold uppercase tracking-wider text-slate-500 mb-2">
              Adresse Email
            </label>
            <input
              id="email"
              v-model="email"
              type="email"
              autocomplete="email"
              required
              placeholder="nom@entreprise.com"
              class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm placeholder-slate-400 focus:outline-none focus:border-blue-600 focus:ring-4 focus:ring-blue-50 transition-all duration-150 bg-slate-50/50 focus:bg-white"
            />
          </div>

          <!-- Champ Mot de passe -->
          <div>
            <div class="flex justify-between items-center mb-2">
              <label for="password" class="block text-xs font-semibold uppercase tracking-wider text-slate-500">
                Mot de passe
              </label>
              <!-- Ajout d'un lien secondaire élégant, classique sur les interfaces pro -->
              <a href="#" class="text-xs font-medium text-blue-600 hover:text-blue-700 transition-colors">
                Oublié ?
              </a>
            </div>
            <input
              id="password"
              v-model="password"
              type="password"
              autocomplete="current-password"
              required
              placeholder="••••••••"
              class="w-full border border-slate-200 rounded-xl px-4 py-3 text-sm placeholder-slate-400 focus:outline-none focus:border-blue-600 focus:ring-4 focus:ring-blue-50 transition-all duration-150 bg-slate-50/50 focus:bg-white"
            />
          </div>

          <!-- Message d'erreur élégant -->
          <div
            v-if="error"
            class="p-3 bg-rose-50 border border-rose-100 text-rose-700 rounded-xl text-xs flex items-center gap-2 animate-fadeIn"
          >
            <span class="w-1.5 h-1.5 rounded-full bg-rose-500 flex-shrink-0"></span>
            <span>{{ error }}</span>
          </div>

          <!-- Bouton de soumission avec état Loading -->
          <button
            type="submit"
            :disabled="isLoading"
            class="w-full bg-slate-900 hover:bg-slate-800 text-white font-medium rounded-xl py-3 text-sm transition-all duration-150 shadow-sm shadow-slate-900/10 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="isLoading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
            <span>{{ isLoading ? 'Connexion en cours...' : 'Se connecter' }}</span>
          </button>

        </form>
      </div>

      <!-- Footer discret -->
      <p class="text-center text-xs text-slate-500 mt-8">
        &copy; 2026 SecureAsset Manager. Tous droits réservés.
      </p>
    </div>

  </div>
</template>

<style scoped>
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-4px); }
  to { opacity: 1; transform: translateY(0); }
}
.animate-fadeIn {
  animation: fadeIn 0.2s ease-out forwards;
}
</style>