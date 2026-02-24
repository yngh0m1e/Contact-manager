<template>
  <div class="min-h-screen bg-[#0f172a] text-white">
    <div class="max-w-6xl mx-auto p-8">
      <div class="flex justify-between items-center mb-12">
        <h1 class="text-5xl font-light">Мои контакты</h1>
        <div class="flex gap-4">
          <button @click="router.push('/profile')" class="px-8 py-3 rounded-2xl border border-white/30 hover:bg-white/10">👤 Профиль</button>
          <button @click="logout" class="px-8 py-3 rounded-2xl bg-red-600/80 hover:bg-red-700">Выйти</button>
        </div>
      </div>

      <div class="flex gap-4 mb-10">
        <input v-model="search" @input="filterContacts"
          placeholder="Поиск..." 
          class="flex-1 bg-white/10 border border-white/20 rounded-3xl px-8 py-5 text-lg focus:border-[#9f1239]" />
        <button @click="openAddModal" class="bg-[#9f1239] hover:bg-[#b91c1c] px-10 rounded-3xl font-medium">+ Новый контакт</button>
      </div>

      <div class="bg-white/5 border border-white/10 rounded-3xl overflow-hidden">
        <table class="w-full">
          <thead>
            <tr class="border-b border-white/10">
              <th class="text-left py-7 px-10 text-zinc-400">Имя</th>
              <th class="text-left py-7 px-10 text-zinc-400">Email</th>
              <th class="text-left py-7 px-10 text-zinc-400">Телефон</th>
              <th class="w-32"></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="c in filteredContacts" :key="c.id" class="border-b border-white/10 hover:bg-white/5">
              <td class="py-7 px-10">{{ c.name }}</td>
              <td class="py-7 px-10 text-zinc-400">{{ c.email }}</td>
              <td class="py-7 px-10 text-zinc-400">{{ c.phone }}</td>
              <td class="py-7 px-10 text-right space-x-6">
                <button @click="editContact(c)" class="text-2xl hover:text-white">✏️</button>
                <button @click="deleteContact(c)" class="text-2xl hover:text-red-400">🗑</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Модальное окно -->
    <dialog ref="modal" class="modal">
      <div class="modal-box bg-[#1a2338] border border-white/10 rounded-3xl">
        <h3 class="text-3xl font-light mb-8">{{ isEdit ? 'Редактировать' : 'Новый контакт' }}</h3>
        <div class="space-y-6">
          <input v-model="form.name" placeholder="Имя" class="w-full bg-white/10 border border-white/20 rounded-2xl px-6 py-4" />
          <input v-model="form.email" type="email" placeholder="Email" class="w-full bg-white/10 border border-white/20 rounded-2xl px-6 py-4" />
          <input v-model="form.phone" placeholder="Телефон" class="w-full bg-white/10 border border-white/20 rounded-2xl px-6 py-4" />
        </div>
        <div class="modal-action">
          <button @click="closeModal" class="btn btn-ghost">Отмена</button>
          <button @click="saveContact" class="bg-[#9f1239] hover:bg-[#b91c1c] px-10 py-3 rounded-2xl">Сохранить</button>
        </div>
      </div>
    </dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useRouter } from 'vue-router'
import axios from 'axios'

const auth = useAuthStore()
const router = useRouter()

const contacts = ref([])
const filteredContacts = ref([])
const search = ref('')
const modal = ref(null)
const isEdit = ref(false)
const editingId = ref(null)
const form = ref({ name: '', email: '', phone: '' })

const loadContacts = async () => {
  const res = await axios.get('/contacts')
  contacts.value = res.data
  filteredContacts.value = res.data
}

const filterContacts = () => {
  const term = search.value.toLowerCase()
  filteredContacts.value = contacts.value.filter(c => 
    c.name.toLowerCase().includes(term) || c.email.toLowerCase().includes(term) || c.phone.toLowerCase().includes(term)
  )
}

const openAddModal = () => { isEdit.value = false; form.value = { name:'', email:'', phone:'' }; modal.value.showModal() }
const editContact = (c) => { isEdit.value = true; editingId.value = c.id; form.value = {...c}; modal.value.showModal() }
const saveContact = async () => {
  if (isEdit.value) await axios.put(`/contacts/${editingId.value}`, form.value)
  else await axios.post('/contacts', form.value)
  closeModal()
  loadContacts()
}
const deleteContact = async (c) => {
  if (confirm(`Удалить ${c.name}?`)) {
    await axios.delete(`/contacts/${c.id}`)
    loadContacts()
  }
}
const closeModal = () => modal.value.close()
const logout = () => auth.logout()

onMounted(() => {
  if (!localStorage.getItem('token')) router.push('/')
  loadContacts()
})
</script>