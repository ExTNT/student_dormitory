import { http } from './http';
import type { AttachmentMeta } from '@/types';

export const attachmentApi = {
  upload(data: { file: File; owner_type: string; owner_id: number; category: string; sort_order?: number }) {
    const form = new FormData();
    form.append('file', data.file);
    form.append('owner_type', data.owner_type);
    form.append('owner_id', String(data.owner_id));
    form.append('category', data.category);
    if (data.sort_order !== undefined) form.append('sort_order', String(data.sort_order));
    return http.post<AttachmentMeta>('/attachments', form).then((res) => res.data);
  },
  list(params: { owner_type: string; owner_id: number; category?: string }) {
    return http.get<AttachmentMeta[]>('/attachments', { params }).then((res) => res.data);
  },
  blob(id: number) {
    return http.get<Blob>(`/attachments/${id}`, { responseType: 'blob' }).then((res) => res.data);
  },
};
