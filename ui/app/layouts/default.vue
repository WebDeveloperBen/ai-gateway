<script lang="ts" setup>
import {
  Bot,
  Command,
  GalleryVerticalEnd,
  Settings2,
  ChevronRight,
  ChevronsUpDown,
  Plus,
  Users,
  Layers,
  Shield,
  FileText,
  Server,
  TriangleAlert,
  LineChart,
  Activity,
  Zap,
  TestTube,
  Globe
} from "lucide-vue-next"

// Dynamic breadcrumbs
const { breadcrumbItems } = useBreadcrumbs()
const {
  app: { name }
} = useAppConfig()
// This is sample data.
const data = {
  user: {
    name: "breezy",
    email: "m@example.com",
    avatar: "https://behonbaker.com/icon.png"
  },
  teams: [
    {
      name,
      logo: GalleryVerticalEnd,
      plan: "Production"
    },
    {
      name: "Development",
      logo: Command,
      plan: "Dev"
    }
  ],
  environments: [
    {
      name: "Production",
      logo: Zap,
      description: "Live environment",
      status: "active"
    },
    {
      name: "Staging",
      logo: TestTube,
      description: "Pre-production testing",
      status: "active"
    },
    {
      name: "Development",
      logo: Command,
      description: "Development environment",
      status: "active"
    },
    {
      name: "Testing",
      logo: Globe,
      description: "QA testing environment",
      status: "maintenance"
    }
  ],
  navMain: [
    {
      title: "Applications",
      url: "/applications",
      icon: Layers,
      isActive: true
    },
    {
      title: "Environments",
      url: "/environments",
      icon: Server
    },
    {
      title: "Models",
      url: "/models",
      icon: Bot
    }
  ],
  navPlayground: [
    {
      title: "Prompts",
      icon: Layers,
      url: "/prompts",
      isActive: true
    },
    {
      title: "Playground",
      url: "/playground",
      icon: Server
    }
  ],
  navObservability: [
    { title: "Analytics", url: "/observability/analytics", icon: LineChart },
    {
      title: "Monitoring",
      icon: Activity,
      isActive: true,
      items: [
        { title: "Logs", url: "/observability/logs", icon: FileText },
        { title: "Alerts", url: "/observability/alerts", icon: TriangleAlert }
      ]
    }
  ],
  navGovernance: [
    {
      title: "Policies",
      url: "/governance/policy",
      icon: Shield
    },
    {
      title: "Audit",
      url: "/governance/audit",
      icon: FileText
    }
  ],
  navAdmin: [
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
      title: "Settings",
      url: "/settings",
      icon: Settings2,
      items: [
        {
          title: "General",
          url: "/settings"
        },
        {
          title: "Security",
          url: "/settings/security"
        }
      ]
    }
  ]
}
const activeTeam = ref(data.teams[0]!)
const activeEnvironment = ref(data.environments[0]!)

// Environment context
const { setEnvironment } = useEnvironment()

// Set initial environment
onMounted(() => {
  setEnvironment(activeEnvironment.value)
})

// Watch for environment changes and update global state
watch(activeEnvironment, (newEnv) => {
  setEnvironment(newEnv)
})

useSeoMeta({ title: "LLM Gateway - Admin Dashboard" })
</script>
<template>
  <UiSidebarProvider v-slot="{ isMobile }">
    <!-- App Sidebar -->
    <UiSidebar collapsible="icon">
      <UiSidebarHeader>
        <UiSidebarMenu>
          <!-- Environment Selector -->
          <UiSidebarMenuItem>
            <UiDropdownMenu>
              <UiDropdownMenuTrigger as-child>
                <UiSidebarMenuButton
                  size="lg"
                  class="group-data-[collapsible=icon]:!size-8 group-data-[collapsible=icon]:!p-0 data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
                >
                  <div
                    class="flex aspect-square size-8 items-center justify-center rounded-lg bg-blue-100 text-blue-600 dark:bg-blue-900/20 dark:text-blue-400"
                  >
                    <component :is="activeEnvironment.logo" class="size-4" />
                  </div>
                  <div class="grid flex-1 text-left text-sm leading-tight">
                    <span class="truncate font-semibold">
                      {{ activeEnvironment.name }}
                    </span>
                    <span class="truncate text-xs">{{ activeEnvironment.description }}</span>
                  </div>
                  <component :is="ChevronsUpDown" class="ml-auto" />
                </UiSidebarMenuButton>
              </UiDropdownMenuTrigger>
              <UiDropdownMenuContent
                class="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
                align="start"
                :side="isMobile ? 'bottom' : 'right'"
                :side-offset="4"
              >
                <UiDropdownMenuLabel class="text-xs text-muted-foreground">Environments</UiDropdownMenuLabel>
                <template v-for="(environment, index) in data.environments" :key="index">
                  <UiDropdownMenuItem
                    class="cursor-pointer gap-2 p-2"
                    :class="[environment.name == activeEnvironment.name && 'bg-muted']"
                    @click="activeEnvironment = environment"
                  >
                    <div class="flex size-6 items-center justify-center rounded-sm border">
                      <component :is="environment.logo" class="size-4 shrink-0" />
                    </div>
                    <div class="flex-1">
                      <div class="font-medium">{{ environment.name }}</div>
                      <div class="text-xs text-muted-foreground">{{ environment.description }}</div>
                    </div>
                    <div
                      :class="[
                        'h-2 w-2 rounded-full',
                        environment.status === 'active' ? 'bg-green-500' : 'bg-orange-500'
                      ]"
                    />
                    <UiDropdownMenuShortcut>⌘{{ index + 1 }}</UiDropdownMenuShortcut>
                  </UiDropdownMenuItem>
                </template>
                <UiDropdownMenuSeparator />
                <UiDropdownMenuItem class="gap-2 p-2">
                  <div class="flex size-6 items-center justify-center rounded-md border bg-background">
                    <component :is="Plus" class="size-4" />
                  </div>
                  <div class="font-medium text-muted-foreground">Add environment</div>
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
            <UiSidebarMenuItem v-for="(item, index) in data.navMain" :key="index">
              <UiSidebarMenuButton as-child :tooltip="item.title">
                <NuxtLink :href="item.url">
                  <component :is="item.icon" />
                  <span>{{ item.title }}</span>
                </NuxtLink>
              </UiSidebarMenuButton>
            </UiSidebarMenuItem>
          </UiSidebarMenu>
        </UiSidebarGroup>

        <!--Engineering -->
        <UiSidebarGroup>
          <UiSidebarGroupLabel label="Engineering" />
          <UiSidebarMenu>
            <template v-for="(item, index) in data.navPlayground" :key="index">
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
        <!-- Observability -->
        <UiSidebarGroup>
          <UiSidebarGroupLabel label="Observability" />
          <UiSidebarMenu>
            <template v-for="(item, index) in data.navObservability" :key="index">
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
        <UiSidebarGroup>
          <UiSidebarGroupLabel label="Administration" />
          <UiSidebarMenu>
            <template v-for="(item, index) in data.navAdmin" :key="index">
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
      <UiNavbar class="flex relative h-16 shrink-0 items-center gap-2 border-b px-4">
        <UiSidebarTrigger class="-ml-1" />
        <UiDivider orientation="vertical" class="mr-2 h-4 w-px" />
        <UiBreadcrumbs :items="breadcrumbItems" />
        <div class="ml-auto">
          <ThemeToggle />
        </div>
      </UiNavbar>
      <UiContainer>
        <slot />
      </UiContainer>
    </UiSidebarInset>
  </UiSidebarProvider>
</template>
