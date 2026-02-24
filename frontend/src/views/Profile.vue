<template>
  <div class="min-h-screen bg-[#0f172a] text-white flex items-center justify-center">
    <div class="w-full max-w-md px-6">
      <div class="text-center mb-12">
        <div class="mx-auto w-24 h-24 bg-white/10 backdrop-blur-3xl rounded-3xl flex items-center justify-center text-7xl mb-6 border border-white/30">
          👤
        </div>
        <h1 class="text-4xl font-light">Мой профиль</h1>
      </div>

      <div class="bg-white/10 backdrop-blur-3xl border border-white/20 rounded-3xl p-10">
        <div class="mb-10 text-center">
          <div class="text-3xl font-light">{{ user.name }}</div>
          <div class="text-zinc-400">{{ user.email }}</div>
        </div>

        <div class="space-y-6">
          <div>
            <label class="block text-zinc-400 mb-2">Имя</label>
            <input v-model="editForm.name" class="w-full bg-white/10 border border-white/20 rounded-2xl px-6 py-4 text-white" />
          </div>
          <div>
            <label class="block text-zinc-400 mb-2">Email</label>
            <input v-model="editForm.email" type="email" class="w-full bg-white/10 border border-white/20 rounded-2xl px-6 py-4 text-white" />
          </div>
        </div>

        <button @click="saveProfile" class="mt-10 w-full bg-[#9f1239] hover:bg-[#b91c1c] py-4 rounded-2xl font-medium transition">
          Сохранить изменения
        </button>
      </div>

      <div class="flex gap-4 justify-center mt-8">
        <button @click="router.push('/contacts')" class="px-8 py-3 rounded-2xl border border-white/20 hover:bg-white/10">📇 Контакты</button>
        <button @click="logout" class="px-8 py-3 rounded-2xl bg-red-700 hover:bg-red-800">Выйти</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import axios from 'axios'

const auth = useAuthStore()
const router = useRouter()

const user = ref({ name: '', email: '' })
const editForm = ref({ name: '', email: '' })

const loadProfile = async () => {
  const res = await axios.get('/profile')
  user.value = res.data
  editForm.value = { name: res.data.name, email: res.data.email }
}

const saveProfile = async () => {
  await axios.put('/profile', editForm.value)
  alert('Профиль обновлён!')
  loadProfile()
}

const logout = () => auth.logout()

onMounted(() => {
  if (!localStorage.getItem('token')) router.push('/')
  loadProfile()
})
</script>