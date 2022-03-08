<template>
  <v-container
    id="dashboard-view"
    fluid
    tag="section"
  >
    <v-row>
      Topic: {{ topic }}<br>
      Message: {{ msg }}
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-row>
          <v-col
            v-for="(emitter, i) in emitters"
            :key="`emitter-${i}`"
            cols="12"
            md="6"
            lg="4"
          >
            <material-emitter-card :emitter="emitter" />
          </v-col>
        </v-row>
      </v-col>

      <v-col
        v-for="({ actionIcon, actionText, ...attrs }, i) in stats"
        :key="i"
        cols="12"
        md="6"
        lg="3"
      >
        <material-stat-card v-bind="attrs">
          <template #actions>
            <v-icon
              class="mr-2"
              small
              v-text="actionIcon"
            />
            <div class="text-truncate">
              {{ actionText }}
            </div>
          </template>
        </material-stat-card>
      </v-col>
      <error-logs-tables-view />
    </v-row>
  </v-container>
</template>

<script>
  // Utilities
  import ErrorLogsTablesView from './ErrorLogs'
  import WebsocketService from '@/services/websocket.service'
  // import ApiService from '@/services/api.service'

  export default {
    name: 'DashboardView',

    components: {
      ErrorLogsTablesView,
    },

    data: () => ({
      stats: [],
      emitters: [],
      topic: '',
      msg: '',
      tabs: 0,
    }),

    computed: {
    },

    created () {
      var t = this
      WebsocketService.topic('file', function (topic, message) {
        t.topic = topic
        t.msg += new Date().toString()
      })
    },
  }
</script>
