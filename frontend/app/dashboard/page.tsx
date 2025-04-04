'use client';

import { useEffect, useState } from 'react';
import { useAuth } from '@/context/auth-context';
import { fetchApi, ApiError } from '@/lib/api';
import ProtectedRoute from '@/components/auth/protected-route';
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar';
import { AppSidebar } from '@/components/dashboard/dashboard-sidebar/dashboard-sidebar';
import { PageHeader } from '@/components/dashboard/page-header';
// import Header from '@/components/dashboard/Header';

interface UserData {
  id: string;
  email: string;
  name: string;
  avatarUrl?: string;
  role?: string;
  createdAt?: string;
}

export default function DashboardPage() {

  return (
    <ProtectedRoute>
        <SidebarProvider>
            <AppSidebar variant="inset" />
            <SidebarInset>
                <PageHeader />
                {/*
                <div className="flex flex-1 flex-col">
                <div className="@container/main flex flex-1 flex-col gap-2">
                    <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
                    <SectionCards />
                    <div className="px-4 lg:px-6">
                        <ChartAreaInteractive />
                    </div>
                    <DataTable data={data} />
                    </div>
                </div>
                </div> */}
        </SidebarInset>
        </SidebarProvider>
    </ProtectedRoute>
  );
}