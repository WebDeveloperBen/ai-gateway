<script lang="ts" setup>
import {
  Bot,
  Command,
  GalleryVerticalEnd,
  Settings2,
  ChevronRight,
  ChevronsUpDown,
  Plus,
  BarChart3,
  Users,
  Layers,
  Shield,
  FileText,
  Lock,
  Globe,
  Bell,
  Server
} from "lucide-vue-next"

// Breadcrumb items
const breadcrumbItems = [
  { label: "LLM Gateway", link: "#" },
  { label: "Dashboard", link: "#" }
]

// This is sample data.
const data = {
  user: {
    name: "breezy",
    email: "m@example.com",
    avatar: "https://behonbaker.com/icon.png"
  },
  teams: [
    {
      name: "LLM Gateway",
      logo: GalleryVerticalEnd,
      plan: "Production"
    },
    {
      name: "Development",
      logo: Command,
      plan: "Dev"
    }
  ],
  navMain: [
    {
      title: "Applications",
      url: "/applications",
      icon: Layers,
      isActive: true,
      items: [
        {
          title: "Overview",
          url: "/applications"
        },
        {
          title: "API Keys",
          url: "/applications/keys"
        }
      ]
    },
    {
      title: "Models",
      url: "/models",
      icon: Bot,
      items: [
        {
          title: "Deployments",
          url: "/models/deployments"
        },
        {
          title: "Providers",
          url: "/models/providers"
        }
      ]
    },
    {
      title: "Users",
      url: "/users",
      icon: Users,
      items: [
        {
          title: "Overview",
          url: "/users"
        },
        {
          title: "Teams",
          url: "/users/teams"
        },
        {
          title: "Roles",
          url: "/users/roles"
        }
      ]
    },
    {
      title: "Environments",
      url: "/environments",
      icon: Server,
      items: [
        {
          title: "Overview",
          url: "/environments"
        },
        {
          title: "Configuration",
          url: "/environments/config"
        }
      ]
    }
  ],
  navObservability: [
    {
      title: "Metrics",
      url: "/metrics",
      icon: Shield
    },
    {
      title: "Audit Logs",
      url: "/audit-logs",
      icon: FileText
    }
  ],
  navGovernance: [
    {
      title: "Policies",
      url: "/policies",
      icon: Shield
    },
    {
      title: "Audit Logs",
      url: "/audit-logs",
      icon: FileText
    }
  ],
  navSettings: [
    {
      title: "Settings",
      url: "/settings",
      icon: Settings2,
      items: [
        {
          title: "General",
          url: "/settings/general"
        },
        {
          title: "Security",
          url: "/settings/security"
        }
      ]
    }
  ],
  projects: []
}
const activeTeam = ref(data.teams[0]!)
useSeoMeta({ title: "LLM Gateway - Admin Dashboard" })
</script>
<template>
  <UiSidebarProvider v-slot="{ isMobile }">
    <!-- App Sidebar -->
    <UiSidebar collapsible="icon">
      <!-- Team switcher -->
      <UiSidebarHeader>
        <UiSidebarMenu>
          <UiSidebarMenuItem>
            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiSidebarMenuButton
                  size="lg"
                  class="group-data-[collapsible=icon]:!size-8 group-data-[collapsible=icon]:!p-0 data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <div
                    class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground"
                  >
                    <component :is="activeTeam.logo" class="size-4" />
                  </div>
                  <div class="grid flex-1 text-left text-sm leading-tight">
                    <span class="truncate font-semibold">
                      {{ activeTeam.name }}
                    </span>
                    <span class="truncate text-xs">{{ activeTeam.plan }}</span>
                  </div>
                  <component :is="" class="ml-auto" />
                  <component :is="ChevronsUpDown" class="ml-auto" />
                </UiSidebarMenuButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent
                class="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                align="start"
                :side="isMobile ? 'bottom' : 'right'"
                :side-offset="4"
              >
                <UiDropdownMenuLabel class="text-xs text-muted-foreground"> Teams </UiDropdownMenuLabel>
                <template v-for="(team, index) in data.teams" :key="index">
                  <UiDropdownMenuItem
                    class="cursor-pointer gap-2 p-2"
                    :class="[team.name == activeTeam.name && 'bg-muted']"
                    @click="activeTeam = team"
                  >
                    <div class="flex size-6 items-center justify-center rounded-sm border">
                      <component :is="team.logo" class="size-4 shrink-0" />
                    </div>
                    {{ team.name }}
                    <UiDropdownMenuShortcut>⌘{{ index + 1 }}</UiDropdownMenuShortcut>
                  </UiDropdownMenuItem>
                </template>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="gap-2 p-2">
                  <div class="flex size-6 items-center justify-center rounded-md border bg-background">
                    <component :is="Plus" class="size-4" />
                  </div>
                  <div class="font-medium text-muted-foreground">Add team</div>
                </UiDropdownMenuItem>
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </UiSidebarMenuItem>
        </UiSidebarMenu>
      </UiSidebarHeader>
      <UiSidebarContent>
        <!-- Platform -->
        <UiSidebarGroup>
          <UiSidebarGroupLabel label="Platform" />
          <UiSidebarMenu>
            <template v-for="(item, index) in data.navMain" :key="index">
              <!-- Items with sub-items -->
              <UiCollapsible v-if="item.items" v-slot="{ open }" as-child :default-open="item.isActive">
                <UiSidebarMenuItem>
                  <UiCollapsibleTrigger as-child>
                    <UiSidebarMenuButton :tooltip="item.title">
                      <component :is="item.icon" />
                      <span>{{ item.title }}</span>
                      <component
                        :is="ChevronRight"
                        class="ml-auto transition-transform duration-200"
                        :class="[open && 'rotate-90']"
                      />
                    </UiSidebarMenuButton>
                  </UiCollapsibleTrigger>
                  <UiCollapsibleContent>
                    <UiSidebarMenuSub>
                      <UiSidebarMenuSubItem v-for="subItem in item.items" :key="subItem.title">
                        <UiSidebarMenuSubButton as-child>
                          <NuxtLink :href="subItem.url">
                            <span>{{ subItem.title }}</span>
                          </NuxtLink>
                        </UiSidebarMenuSubButton>
                      </UiSidebarMenuSubItem>
                    </UiSidebarMenuSub>
                  </UiCollapsibleContent>
                </UiSidebarMenuItem>
              </UiCollapsible>
              <!-- Items without sub-items -->
              <UiSidebarMenuItem v-else>
                <UiSidebarMenuButton as-child :tooltip="item.title">
                  <NuxtLink :href="item.url">
                    <component :is="item.icon" />
                    <span>{{ item.title }}</span>
                  </NuxtLink>
                </UiSidebarMenuButton>
              </UiSidebarMenuItem>
            </template>
          </UiSidebarMenu>
        </UiSidebarGroup>

        <!-- Observability -->
        <UiSidebarGroup>
          <UiSidebarGroupLabel label="Observability" />
          <UiSidebarMenu>
            <UiSidebarMenuItem v-for="(item, index) in data.navObservability" :key="index">
              <UiSidebarMenuButton as-child :tooltip="item.title">
                <NuxtLink :href="item.url">
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </NuxtLink>
              </UiSidebarMenuButton>
            </UiSidebarMenuItem>
          </UiSidebarMenu>
        </UiSidebarGroup>
        <!-- Governance -->
        <UiSidebarGroup>
          <UiSidebarGroupLabel label="Governance" />
          <UiSidebarMenu>
            <UiSidebarMenuItem v-for="(item, index) in data.navGovernance" :key="index">
              <UiSidebarMenuButton as-child :tooltip="item.title">
                <NuxtLink :href="item.url">
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </NuxtLink>
              </UiSidebarMenuButton>
            </UiSidebarMenuItem>
          </UiSidebarMenu>
        </UiSidebarGroup>
        <!-- Settings -->
        <UiSidebarGroup>
          <UiSidebarGroupLabel label="Settings" />
          <UiSidebarMenu>
            <template v-for="(item, index) in data.navSettings" :key="index">
              <!-- Items with sub-items -->
              <UiCollapsible v-if="item.items" v-slot="{ open }" as-child>
                <UiSidebarMenuItem>
                  <UiCollapsibleTrigger as-child>
                    <UiSidebarMenuButton :tooltip="item.title">
                      <component :is="item.icon" />
                      <span>{{ item.title }}</span>
                      <component
                        :is="ChevronRight"
                        class="ml-auto transition-transform duration-200"
                        :class="[open && 'rotate-90']"
                      />
                    </UiSidebarMenuButton>
                  </UiCollapsibleTrigger>
                  <UiCollapsibleContent>
                    <UiSidebarMenuSub>
                      <UiSidebarMenuSubItem v-for="subItem in item.items" :key="subItem.title">
                        <UiSidebarMenuSubButton as-child>
                          <NuxtLink :href="subItem.url">
                            <span>{{ subItem.title }}</span>
                          </NuxtLink>
                        </UiSidebarMenuSubButton>
                      </UiSidebarMenuSubItem>
                    </UiSidebarMenuSub>
                  </UiCollapsibleContent>
                </UiSidebarMenuItem>
              </UiCollapsible>
              <!-- Items without sub-items -->
              <UiSidebarMenuItem v-else>
                <UiSidebarMenuButton as-child :tooltip="item.title">
                  <NuxtLink :href="item.url">
                    <component :is="item.icon" />
                    <span>{{ item.title }}</span>
                  </NuxtLink>
                </UiSidebarMenuButton>
              </UiSidebarMenuItem>
            </template>
          </UiSidebarMenu>
        </UiSidebarGroup>
      </UiSidebarContent>
      <UiSidebarRail />
      <!-- Footer-->
      <UiSidebarFooter>
        <UiSidebarMenu>
          <UiSidebarMenuItem>
            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiSidebarMenuButton
                  size="lg"
                  class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <UiAvatar class="size-8 rounded-lg">
                    <UiAvatarImage :src="data.user.avatar" :alt="data.user.name" />
                    <UiAvatarFallback class="rounded-lg">BB</UiAvatarFallback>
                  </UiAvatar>
                  <div class="grid flex-1 text-left text-sm leading-tight">
                    <span class="truncate font-semibold">{{ data.user.name }}</span>
                    <span class="truncate text-xs">{{ data.user.email }}</span>
                  </div>
                  <component :is="ChevronsUpDown" class="ml-auto size-4" />
                </UiSidebarMenuButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent
                class="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                :side="isMobile ? 'bottom' : 'right'"
                :side-offset="4"
                align="end"
              >
                <UiDropdownMenuLabel class="p-0 font-normal">
                  <div class="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
                    <UiAvatar class="size-8 rounded-lg">
                      <UiAvatarImage :src="data.user.avatar" :alt="data.user.name" />
                      <UiAvatarFallback class="rounded-lg">BB</UiAvatarFallback>
                    </UiAvatar>
                    <div class="grid flex-1 text-left text-sm leading-tight">
                      <span class="truncate font-semibold">{{ data.user.name }}</span>
                      <span class="truncate text-xs">{{ data.user.email }}</span>
                    </div>
                  </div>
                </UiDropdownMenuLabel>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuGroup>
                  <UiDropdownMenuItem icon="Sparkles," title="Upgrade to Pro" />
                </UiDropdownMenuGroup>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuGroup>
                  <UiDropdownMenuItem icon="lucide:badge-check" title="Account" />
                  <UiDropdownMenuItem icon="lucide:credit-card" title="Billing" />
                  <UiDropdownMenuItem icon="lucide:settings-2" title="Settings" />
                  <UiDropdownMenuItem icon="lucide:bell" title="Notifications" />
                </UiDropdownMenuGroup>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem icon="lucide:log-out" title="Log out" />
              </UiDropdownMenuContent>
            </UiDropdownMenu>
          </UiSidebarMenuItem>
        </UiSidebarMenu>
      </UiSidebarFooter>
    </UiSidebar>
    <!-- Sidebar main content -->
    <UiSidebarInset>
      <!-- Navbar -->
      <UiNavbar sticky class="flex relative h-16 shrink-0 items-center gap-2 border-b px-4">
        <UiSidebarTrigger class="-ml-1" />
        <UiDivider orientation="vertical" class="mr-2 h-4 w-px" />
        <UiBreadcrumbs :items="breadcrumbItems" />
        <ThemeToggle />
      </UiNavbar>
      <PageContainer>
        <slot />
      </PageContainer>
    </UiSidebarInset>
  </UiSidebarProvider>
</template>
