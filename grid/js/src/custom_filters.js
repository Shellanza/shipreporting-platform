export function CustomDateComponent () {
}

CustomDateComponent.prototype.init = function (params) {
  var template =
        "<input type='text' data-input />" +
        "<a class='input-button' title='clear' data-clear>" +
            "<i class='fa fa-times'></i>" +
        '</a>'

  this.params = params

  this.eGui = document.createElement('div')

  var eGui = this.eGui

  eGui.setAttribute('role', 'presentation')
  eGui.classList.add('ag-input-wrapper')
  eGui.classList.add('custom-date-filter')
  eGui.innerHTML = template

  this.eInput = eGui.querySelector('input')

  // eslint-disable-next-line no-undef
  this.picker = flatpickr(this.eGui, {
    onChange: this.onDateChanged.bind(this),
    dateFormat: 'd/m/Y',
    wrap: true
  })

  this.picker.calendarContainer.classList.add('ag-custom-component-popup')

  this.date = null
}

CustomDateComponent.prototype.getGui = function () {
  return this.eGui
}

CustomDateComponent.prototype.onDateChanged = function (selectedDates) {
  this.date = selectedDates[0] || null
  this.params.onDateChanged()
}

CustomDateComponent.prototype.getDate = function () {
  return this.date
}

CustomDateComponent.prototype.setDate = function (date) {
  this.picker.setDate(date)
  this.date = date
}
