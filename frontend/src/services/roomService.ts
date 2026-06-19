import axios from 'axios'

export interface Room {
  id: number
  name: string
  price: number
  image: string
  capacity: number
  description: string
  features: string[]
}

export const roomService = {
  getRooms: async (): Promise<Room[]> => {
    try {
      const response = await axios.get('/api/rooms')
      if (response.data && response.data.success) {
        return response.data.data
      }
      return []
    } catch (error) {
      console.error('Error fetching rooms:', error)
      return []
    }
  },
  getRoomById: async (id: number): Promise<Room | undefined> => {
    try {
      const response = await axios.get(`/api/rooms/${id}`)
      if (response.data && response.data.success) {
        return response.data.data
      }
      return undefined
    } catch (error) {
      console.error('Error fetching room details:', error)
      return undefined
    }
  }
}
