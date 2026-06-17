export interface Room {
  id: number
  name: string
  price: number
  image: string
  capacity: number
  description: string
  features: string[]
}

const mockRooms: Room[] = [
  {
    id: 1,
    name: 'Standard Room',
    price: 350000,
    image: 'https://images.unsplash.com/photo-1566665797739-1674de7a421a?w=600&auto=format&fit=crop&q=60',
    capacity: 2,
    description: 'A cozy standard room featuring a queen size bed, free high speed Wi-Fi, and standard amenities for a comfortable stay.',
    features: ['Queen size bed', 'Air Conditioning', 'Free Wi-Fi', 'TV', 'En-suite bathroom']
  },
  {
    id: 2,
    name: 'Deluxe Room',
    price: 550000,
    image: 'https://images.unsplash.com/photo-1582719508461-905c673771fd?w=600&auto=format&fit=crop&q=60',
    capacity: 2,
    description: 'Spacious deluxe room with premium city views, writing desk, smart TV, mini-fridge, and high-end bath amenities.',
    features: ['King size bed', 'Air Conditioning', 'High Speed Wi-Fi', 'Smart TV', 'Mini Bar', 'City View']
  },
  {
    id: 3,
    name: 'Family Suite',
    price: 850000,
    image: 'https://images.unsplash.com/photo-1591088398332-8a7791972843?w=600&auto=format&fit=crop&q=60',
    capacity: 4,
    description: 'Perfect for families. Includes two separate bedrooms, living area, dining table, and fully stocked kitchenette.',
    features: ['2 King size beds', 'Separate Living Room', 'Kitchenette', 'Smart TV', 'Balcony', 'Breakfast Included']
  }
]

export const roomService = {
  getRooms: async (): Promise<Room[]> => {
    // Simulate API delay
    await new Promise((resolve) => setTimeout(resolve, 200))
    return mockRooms
  },
  getRoomById: async (id: number): Promise<Room | undefined> => {
    await new Promise((resolve) => setTimeout(resolve, 100))
    return mockRooms.find((r) => r.id === id)
  }
}
