<script setup >
import { reservationService } from '@/services/admin/reservationService';
import { onMounted, ref } from 'vue';

  const reservations = ref([])
  
  onMounted(async () => {
     try {
      const data = await reservationService.getAll()
      reservations.value = data
     } catch (error) {
      console.log(error)
     }  
  })
</script>

<template>
  <div class="dashboard-view">
    <div class="stats-grid">
      <div class="stat-card">
        <h3 class="stat-title">Active Reservations</h3>
        <p class="stat-value">{{ reservations.length }}</p>
        <span class="stat-trend trend-up">+12% this week</span>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Occupancy Rate</h3>
        <p class="stat-value">86%</p>
        <span class="stat-trend trend-up">+4% vs last month</span>
      </div>
      <div class="stat-card">
        <h3 class="stat-title">Monthly Revenue</h3>
        <p class="stat-value">$18,450</p>
        <span class="stat-trend trend-up">+8% growth</span>
      </div>
    </div>

    <div class="recent-activity">
      <h2>Recent Bookings</h2>
      <div class="table-container">
        <table class="activity-table">
          <thead>
            <tr>
              <th>Guest Name</th>
              <th>Room</th>
              <th>Status</th>
              <th>Check-in</th>
              <th>Check-out</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="res in reservations" :key="res.id">

<td>{{ res.full_name }}</td>
              <td> {{ res.room.name }}</td>
              <td>
                <span class="badge" :class="{
                  'badge-success': res.status === 'approved',
                  'badge-pending': res.status === 'pending',
                  'badge-danger': res.status === 'cancel'
                }">{{ res.status }}</span>
              </td>
              <td>{{ res.checkin_date }}</td>
              <td>{{ res.checkout_date }}</td>
            
             

            </tr>
            <tr>
              
            </tr>
           
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-view {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
}

.stat-card {
  background-color: #ffffff;
  padding: 24px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
}

.stat-title {
  font-size: 0.875rem;
  color: #64748b;
  margin: 0 0 8px 0;
  font-weight: 600;
}

.btn-warnin g{
  background-color: orange;
  color: #ffffff;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #0f172a;
  margin: 0 0 6px 0;
}

.stat-trend {
  font-size: 0.75rem;
  font-weight: 600;
}

.trend-up {
  color: #10b981;
}

.recent-activity {
  background-color: #ffffff;
  padding: 24px;
  border-radius: 16px;
  border: 1px solid rgba(228, 228, 231, 0.8);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.recent-activity h2 {
  font-size: 1.125rem;
  font-weight: 700;
  margin: 0;
  color: #0f172a;
}

.table-container {
  overflow-x: auto;
}

.activity-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.activity-table th {
  padding: 12px 16px;
  border-bottom: 1px solid #e2e8f0;
  color: #64748b;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.activity-table td {
  padding: 16px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 0.875rem;
  color: #334155;
}

.badge {
  padding: 4px 8px;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.badge-success {
  background-color: #d1fae5;
  color: #065f46;
}

.badge-danger {
  background-color: red;
  color: white;
}

.badge-pending {
  background-color: #fef3c7;
  color: #92400e;
}
</style>
