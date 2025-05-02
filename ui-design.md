# UI/UX Design - Multi-Tenant Messaging System

## 1. Layout Dasar
```
+------------------+-----------------------------------+
|    Sidebar      |           Main Content            |
| +-------------+ | +-------------------------------+ |
| | Tenant List | | |        Message List          | |
| |             | | |                              | |
| | - Tenant A  | | | [Search & Filter Bar]        | |
| | - Tenant B  | | |                              | |
| | - Tenant C  | | | +-------------------------+   | |
| |             | | | | Message Item            |   | |
| | [Add Tenant]| | | | - Content              |   | |
| |             | | | | - Timestamp            |   | |
| |             | | | | - Actions              |   | |
| |             | | | +-------------------------+   | |
| |             | | |                              | |
| +-------------+ | | [Load More / Pagination]     | |
|                 | |                              | |
+------------------+--------------------------------+
```

## 2. Komponen Utama

### 2.1 Tenant Management
- **Tenant List**
  - Daftar tenant dengan status aktif/nonaktif
  - Indikator jumlah pesan yang belum dibaca
  - Quick actions (edit, delete)
  - Tombol "Add New Tenant"

- **Tenant Form**
  ```
  +------------------------+
  |    Add/Edit Tenant    |
  +------------------------+
  | Name: [____________]  |
  | Description: [______] |
  | Worker Count: [___]   |
  |                      |
  | [Cancel] [Save]      |
  +------------------------+
  ```

### 2.2 Message Management
- **Message List**
  - Infinite scroll dengan cursor pagination
  - Filter dan pencarian
  - Grouping berdasarkan tanggal
  - Preview pesan dengan format yang sesuai

- **Message Item**
  ```
  +--------------------------------+
  | Message Title                  |
  | [Payload Preview]              |
  | Timestamp                      |
  |                               |
  | [Reply] [Forward] [Delete]    |
  +--------------------------------+
  ```

- **Message Actions**
  - Quick actions untuk setiap pesan
  - Bulk actions untuk multiple pesan
  - Export/download opsi

### 2.3 Monitoring Dashboard
```
+----------------------------------+
|        System Overview           |
+----------------------------------+
| Active Tenants    | Queue Depth  |
| [Graph]           | [Graph]      |
|                   |              |
| Worker Status     | Error Rate   |
| [Status Cards]    | [Graph]      |
+----------------------------------+
```

## 3. Interaksi & Fitur

### 3.1 Tenant Operations
- Drag & drop untuk reorder tenant list
- Click untuk memilih tenant
- Double click untuk edit tenant
- Context menu untuk quick actions
- Real-time status updates

### 3.2 Message Operations
- Infinite scroll untuk message list
- Pull-to-refresh untuk update terbaru
- Swipe actions pada mobile view
- Preview attachment/media
- Copy message content
- Share message

### 3.3 Search & Filter
```
+----------------------------------------+
| Search: [_____________________] [🔍]   |
|                                        |
| Filters:                               |
| Date: [Start] - [End]                  |
| Type: [Dropdown]                       |
| Status: [Multiple Select]              |
+----------------------------------------+
```

## 4. Responsive Design

### 4.1 Desktop View
- Full sidebar visible
- Multi-column message list
- Advanced filtering options
- Keyboard shortcuts

### 4.2 Tablet View
- Collapsible sidebar
- Single column message list
- Simplified filters
- Touch-friendly controls

### 4.3 Mobile View
- Bottom navigation
- Swipe between views
- Simplified actions
- Optimized for touch

## 5. Theme & Styling

### 5.1 Color Scheme
- Primary: #1E88E5 (Blue)
- Secondary: #7CB342 (Green)
- Error: #E53935 (Red)
- Warning: #FFA000 (Amber)
- Background: #F5F5F5
- Text: #212121

### 5.2 Typography
- Heading: Inter
- Body: Roboto
- Monospace: JetBrains Mono (untuk payload code)

### 5.3 Components
- Material Design based
- Consistent padding (16px)
- Rounded corners (8px)
- Subtle shadows
- Clear visual hierarchy

## 6. Feedback & States

### 6.1 Loading States
- Skeleton screens
- Progress indicators
- Placeholder content

### 6.2 Empty States
```
+--------------------------------+
|          No Messages           |
|            (icon)              |
|     Start the conversation     |
|         [New Message]          |
+--------------------------------+
```

### 6.3 Error States
- Clear error messages
- Retry options
- Fallback views

## 7. Accessibility

- ARIA labels
- Keyboard navigation
- High contrast mode
- Screen reader support
- Focus indicators 