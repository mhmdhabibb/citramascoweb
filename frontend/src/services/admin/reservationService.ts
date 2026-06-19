import type { ApiResponse, Reservation } from "@/types";
import api from "../api";

export const reservationService = {
    /**
     * GET /api/reservations
     */
    getAll: async (): Promise<Reservation[]> => {
        const response = await api.get<ApiResponse<Reservation[]>>('/reservations')
        if (!response.data.data) {
            throw new Error(response.data.message || 'Failed to fetch reservations')
        }
        return response.data.data
    },

    /**
     * PATCH /api/reservation/approve/:id
     */
    approve: async (id: string): Promise<void> => {
        const response = await api.patch<ApiResponse<any>>(`/reservation/approve/${id}`)
        if (!response.data.success) {
            throw new Error(response.data.message || 'Failed to approve reservation')
        }
    },

    /**
     * PATCH /api/reservation/reject/:id
     */
    reject: async (id: string): Promise<void> => {
        const response = await api.patch<ApiResponse<any>>(`/reservation/reject/${id}`)
        if (!response.data.success) {
            throw new Error(response.data.message || 'Failed to reject reservation')
        }
    },

    /**
     * PATCH /api/reservation/cancel/:id
     */
    cancel: async (id: string): Promise<void> => {
        const response = await api.patch<ApiResponse<any>>(`/reservation/cancel/${id}`)
        if (!response.data.success) {
            throw new Error(response.data.message || 'Failed to cancel reservation')
        }
    },

    /**
     * POST /api/reservation
     */
    create: async (data: {
        room_id: string;
        full_name: string;
        email: string;
        check_in_date: string;
        check_out_date: string;
        number_of_guest: number;
        is_offer?: boolean;
        offer_code?: string;
    }): Promise<void> => {
        const response = await api.post<ApiResponse<any>>('/reservation', data)
        if (!response.data.success) {
            throw new Error(response.data.message || 'Failed to submit reservation')
        }
    }
}