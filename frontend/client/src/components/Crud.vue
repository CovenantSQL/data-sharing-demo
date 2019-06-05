<template>
  <v-container>
    <div id="app">
      <v-app id="inspire">
        <div>
          <v-toolbar flat color="white">
            <v-toolbar-title>疫苗供应链管理 Demo</v-toolbar-title>
            <v-divider
                class="mx-2"
                inset
                vertical
            ></v-divider>
            <v-spacer></v-spacer>
            <v-dialog v-model="dialog" max-width="500px">
              <template v-slot:activator="{ on }">
                <v-btn color="primary" dark class="mb-2" v-on="on">New Item</v-btn>
              </template>
              <v-card>
                <v-card-title>
                  <span class="headline">{{ formTitle }}</span>
                </v-card-title>

                <v-flex xs12 sm6 md4>

                </v-flex>


                <v-card-text>
                  <v-container grid-list-md>
                    <v-layout wrap>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.serial" label="Serial"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.factory" label="Factory"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-menu
                            ref="menu"
                            v-model="menu"
                            :close-on-content-click="false"
                            :nudge-right="40"
                            :return-value.sync="date"
                            lazy
                            transition="scale-transition"
                            offset-y
                            full-width
                            min-width="290px"
                        >
                          <template v-slot:activator="{ on }">
                            <v-text-field
                                v-model="editedItem.date"
                                label="Date"
                                readonly
                                v-on="on"
                            ></v-text-field>
                          </template>
                          <v-date-picker v-model="editedItem.date" no-title scrollable>
                            <v-spacer></v-spacer>
                            <v-btn flat color="primary" @click="menu = false">Cancel</v-btn>
                            <v-btn flat color="primary" @click="$refs.menu.save(date)">OK</v-btn>
                          </v-date-picker>
                        </v-menu>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.batch" label="Batch"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.carrier" label="Carrier"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.cold_van" label="ColdVan"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.distributor" label="Distributor"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.hospital" label="Hospital"></v-text-field>
                      </v-flex>
                      <v-flex xs12 sm6 md4>
                        <v-text-field v-model="editedItem.patient" label="Patient"></v-text-field>
                      </v-flex>
                    </v-layout>
                  </v-container>
                </v-card-text>

                <v-card-actions>
                  <v-spacer></v-spacer>
                  <v-btn color="blue darken-1" flat @click="close">Cancel</v-btn>
                  <v-btn color="blue darken-1" flat @click="save">Save</v-btn>
                </v-card-actions>
              </v-card>
            </v-dialog>
          </v-toolbar>
          <v-data-table
              :headers="headers"
              :items="items"
              class="elevation-1"
          >
            <v-progress-linear v-slot:progress color="blue" indeterminate></v-progress-linear>
            <template slot="items" slot-scope="props">
              <td>{{ props.item.id }}</td>
              <td class="text-xs-left">{{ props.item.serial }}</td>
              <td class="text-xs-left">{{ props.item.factory }}</td>
              <td class="text-xs-left">{{ props.item.date }}</td>
              <td class="text-xs-left">{{ props.item.batch }}</td>
              <td class="text-xs-left">{{ props.item.carrier }}</td>
              <td class="text-xs-left">{{ props.item.cold_van }}</td>
              <td class="text-xs-left">{{ props.item.distributor }}</td>
              <td class="text-xs-left">{{ props.item.hospital }}</td>
              <td class="text-xs-left">{{ props.item.patient }}</td>
              <td class="justify-center layout px-0">
                <v-icon
                    small
                    class="mr-2"
                    @click="editItem(props.item)"
                >
                  edit
                </v-icon>
                <v-icon
                    small
                    @click="deleteItem(props.item)"
                >
                  delete
                </v-icon>
              </td>
            </template>
            <template v-slot:no-data>
              <v-btn color="primary" @click="initialize">Reset</v-btn>
            </template>
          </v-data-table>
        </div>
      </v-app>
    </div>
  </v-container>
</template>

<script>
  import axios from 'axios';

  export default {
    data: () => ({
      menu: '',
      date: '',
      picker: '',
      formatted: '',
      dialog: false,
      headers: [
        {text: 'ID', value: 'id', align: 'left'},
        {text: 'Serial', value: 'serial', align: 'left'},
        {text: 'Factory', value: 'factory', align: 'left'},
        {text: 'Date', value: 'date'},
        {text: 'Batch', value: 'batch'},
        {text: 'Carrier', value: 'carrier'},
        {text: 'ColdVan', value: 'cold_van'},
        {text: 'Distributor', value: 'distributor'},
        {text: 'Hospital', value: 'hospital'},
        {text: 'Patient', value: 'patient'},
      ],
      items: [],
      editedIndex: -1,
      editedItem: {
        serial: '',
        factory: '',
        date: '',
        batch: '',
        carrier: '',
        cold_van: '',
        distributor: '',
        hospital: '',
        patient: '',
      },
      defaultItem: {
        serial: '',
        factory: '',
        date: '',
        batch: '',
        carrier: '',
        cold_van: '',
        distributor: '',
        hospital: '',
        patient: '',
      }
    }),

    computed: {
      formTitle() {
        return this.editedIndex === -1 ? 'New Item' : 'Edit Item'
      }
    },

    watch: {
      dialog(val) {
        val || this.close()
      },
    },

    created() {
      this.initialize()
    },

    methods: {
      initialize() {
        axios.defaults.headers.common['Authorization'] = localStorage.getItem('token');
        axios.get('/apiv1/cargos')
          .then(resp => {
            this.items = resp.data;
            if (resp.data == null) {
              this.$notify({
                group: 'crud',
                type: 'error',
                title: "Get cargo failed",
                text: "no cargo got",
              });
            }
            this.$log.debug(resp);
          })
          .catch(err => {
            this.$notify({
              group: 'crud',
              type: 'error',
              title: "Get cargo failed",
              text: err.toString(),
            });
            this.$log.error(err);
          });
      },

      editItem(item) {
        this.editedIndex = this.items.indexOf(item);
        this.editedItem = Object.assign({}, item);
        this.$log.debug("edit: ", this.editedIndex, item);
        this.dialog = true;
      },

      deleteItem(item) {
        const index = this.items.indexOf(item);
        confirm('Are you sure you want to delete this item?') && this.items.splice(index, 1);
      },

      close() {
        this.dialog = false;
        setTimeout(() => {
          this.editedItem = Object.assign({}, this.defaultItem);
          this.editedIndex = -1;
        }, 300)
      },

      save() {
        if (this.editedIndex > -1) {
          let diff = {};
          for (let key in this.editedItem) {
            if (!(key in this.items[this.editedIndex]) ||
              this.editedItem[key] !== this.items[this.editedIndex][key]) {
              diff[key] = this.editedItem[key];
            }
          }
          if (Object.keys(diff).length > 0) {
            diff['id'] = this.editedItem.id;
            this.$log.debug(diff)
            axios.put('/apiv1/cargo', diff)
              .then(resp => {
                this.$log.debug(resp)
              })
              .catch(err => {
                this.$notify({
                  group: 'crud',
                  type: 'error',
                  title: "Modify cargo failed",
                  text: err.toString(),
                });
                this.$log.error(err)
              });
          }

          Object.assign(this.items[this.editedIndex], this.editedItem)
        } else {
          this.items.push(this.editedItem)
        }
        this.close()
      }
    }
  }
</script>

<style>
  table.v-table tbody td:first-child, table.v-table tbody td:not(:first-child), table.v-table tbody th:first-child, table.v-table tbody th:not(:first-child), table.v-table thead td:first-child, table.v-table thead td:not(:first-child), table.v-table thead th:first-child, table.v-table thead th:not(:first-child) {
    padding: 0 12px
  }
</style>
