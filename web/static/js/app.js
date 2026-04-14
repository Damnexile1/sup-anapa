document.addEventListener('DOMContentLoaded', function() {
    loadInstructors();
    loadSlots();
});

function loadInstructors() {
    fetch('/api/instructors')
        .then(function(res) { return res.json(); })
        .then(function(instructors) {
            var grid = document.getElementById('instructors-grid');
            if (!grid) return;
            if (instructors.length === 0) {
                grid.innerHTML = '<p class="text-gray-500 text-center col-span-full">Инструкторы скоро появятся</p>';
                return;
            }
            var placeholder = 'https://via.placeholder.com/400x400?text=Instructor';
            instructors.forEach(function(instructor) {
                var photoURL = instructor.Photo || placeholder;
                var div = document.createElement('div');
                div.className = 'bg-white rounded-lg shadow-lg overflow-hidden hover:shadow-xl transition';

                var img = document.createElement('img');
                img.src = photoURL;
                img.alt = instructor.Name;
                img.className = 'w-full h-64 object-cover';
                img.onerror = function() { this.src = placeholder; };
                div.appendChild(img);

                var body = document.createElement('div');
                body.className = 'p-6';

                var name = document.createElement('h3');
                name.className = 'text-2xl font-bold text-gray-800 mb-2';
                name.textContent = instructor.Name;
                body.appendChild(name);

                var phone = document.createElement('p');
                phone.className = 'text-gray-600 mb-3';
                phone.textContent = instructor.Phone;
                body.appendChild(phone);

                var desc = document.createElement('p');
                desc.className = 'text-gray-700';
                desc.textContent = instructor.Description || 'Опытный инструктор SUP';
                body.appendChild(desc);

                div.appendChild(body);
                grid.appendChild(div);
            });
        });
}

var selectedSlotData = null;

function loadSlots() {
    fetch('/api/slots')
        .then(function(res) { return res.json(); })
        .then(function(slots) {
            var container = document.getElementById('slots-container');
            if (!container) return;
            if (slots.length === 0) {
                container.innerHTML = '<p class="text-gray-500">Нет доступных слотов для бронирования</p>';
                return;
            }

            var slotsByDate = {};
            slots.forEach(function(slot) {
                var d = new Date(slot.Date);
                var dateStr = d.toLocaleDateString('ru-RU');
                if (!slotsByDate[dateStr]) slotsByDate[dateStr] = [];
                slotsByDate[dateStr].push(slot);
            });

            var dates = Object.keys(slotsByDate).sort();
            dates.forEach(function(date) {
                var dateDiv = document.createElement('div');
                dateDiv.className = 'mb-4';

                var h3 = document.createElement('h3');
                h3.className = 'text-lg font-semibold mb-2';
                h3.textContent = date;
                dateDiv.appendChild(h3);

                var gridDiv = document.createElement('div');
                gridDiv.className = 'grid grid-cols-1 md:grid-cols-2 gap-3';

                slotsByDate[date].forEach(function(slot) {
                    var start5 = slot.StartTime ? slot.StartTime.substring(0, 5) : '';
                    var end5 = slot.EndTime ? slot.EndTime.substring(0, 5) : '';
                    var isUnavailable = slot.Status === 'pending' || slot.Status === 'booked';
                    var card = document.createElement('div');
                    card.className = 'border rounded-lg p-4 cursor-pointer transition';
                    if (isUnavailable) {
                        card.className += ' border-gray-300 bg-gray-100 opacity-60 cursor-not-allowed';
                        card.title = slot.Status === 'pending' ? 'Слот временно забронирован, ожидает подтверждения' : 'Слот подтверждён';
                    } else {
                        card.className += ' border-gray-300 hover:border-blue-500';
                        card.onclick = (function(s) {
                            return function() {
                                selectSlot(s.ID, date, start5 + ' - ' + end5, s.Price, 'Инструктор #' + s.InstructorID, s.MaxPeople);
                            };
                        })(slot);
                    }

                    if (isUnavailable) {
                        card.style.textDecoration = 'line-through';
                    }

                    var flex = document.createElement('div');
                    flex.className = 'flex justify-between items-center';

                    var left = document.createElement('div');
                    var timeP = document.createElement('p');
                    timeP.className = 'font-semibold';
                    timeP.textContent = start5 + ' - ' + end5;
                    left.appendChild(timeP);
                    var capP = document.createElement('p');
                    capP.className = 'text-sm text-gray-600';
                    capP.textContent = 'До ' + slot.MaxPeople + ' человек';
                    left.appendChild(capP);

                    var right = document.createElement('div');
                    right.className = 'text-right';
                    var priceP = document.createElement('p');
                    priceP.className = 'text-lg font-bold text-blue-600';
                    priceP.textContent = slot.Price + ' \u20bd';
                    right.appendChild(priceP);

                    var statusBadge = document.createElement('span');
                    statusBadge.className = 'text-xs font-semibold ml-2';
                    if (slot.Status === 'pending') {
                        statusBadge.textContent = 'Занят';
                        statusBadge.className += ' text-orange-500';
                    } else if (slot.Status === 'booked') {
                        statusBadge.textContent = 'Занят';
                        statusBadge.className += ' text-red-500';
                    }
                    right.appendChild(priceP);
                    right.appendChild(statusBadge);

                    flex.appendChild(left);
                    flex.appendChild(right);
                    card.appendChild(flex);
                    gridDiv.appendChild(card);
                });

                dateDiv.appendChild(gridDiv);
                container.appendChild(dateDiv);
            });
        });
}

function selectSlot(slotId, date, time, price, instructor, maxPeople) {
    selectedSlotData = { date: date, time: time, price: price, instructor: instructor };
    document.getElementById('selected-slot-id').value = slotId;
    document.getElementById('people-count').max = maxPeople;
    document.getElementById('booking-form-container').classList.remove('hidden');

    var infoDiv = document.getElementById('weather-info');
    infoDiv.innerHTML = '<h3 class="font-semibold mb-2">Детали бронирования</h3>' +
        '<div class="space-y-2">' +
        '<p class="text-gray-600">Дата: <span class="font-semibold">' + date + '</span></p>' +
        '<p class="text-gray-600">Время: <span class="font-semibold">' + time + '</span></p>' +
        '<p class="text-gray-600">Инструктор: <span class="font-semibold">' + instructor + '</span></p>' +
        '<p class="text-gray-600">Цена: <span class="font-semibold text-blue-600">' + price + ' \u20bd</span></p>' +
        '<p class="text-gray-600">Максимум человек: <span class="font-semibold">' + maxPeople + '</span></p>' +
        '</div>';

    document.getElementById('booking-form-container').scrollIntoView({ behavior: 'smooth' });
}

document.addEventListener('DOMContentLoaded', function() {
    var form = document.getElementById('booking-form');
    if (form) {
        form.addEventListener('submit', function(e) {
            e.preventDefault();
            var formData = new FormData(form);
            var data = {
                slot_id: parseInt(formData.get('slot_id')),
                client_name: formData.get('client_name'),
                client_phone: formData.get('client_phone'),
                client_email: formData.get('client_email'),
                people_count: parseInt(formData.get('people_count'))
            };

            var btn = form.querySelector('button[type="submit"]');
            btn.disabled = true;
            btn.textContent = 'Отправка...';

            fetch('/booking', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(data)
            })
            .then(function(res) {
                if (!res.ok) {
                    return res.json().then(function(err) {
                        throw new Error(err || 'Ошибка при бронировании');
                    });
                }
                return res.json();
            })
            .then(function(result) {
                var slotDetails = selectedSlotData ? selectedSlotData.date + ' ' + selectedSlotData.time : 'Слот #' + data.slot_id;
                var resultDiv = document.getElementById('booking-result');
                resultDiv.innerHTML = '<div class="bg-green-100 border border-green-400 text-green-700 px-4 py-3 rounded">' +
                    '<p class="font-semibold">Бронирование отправлено!</p>' +
                    '<p class="text-sm"><strong>Номер бронирования:</strong> #' + result.ID + '</p>' +
                    '<p class="text-sm"><strong>Клиент:</strong> ' + data.client_name + '</p>' +
                    '<p class="text-sm"><strong>Телефон:</strong> ' + data.client_phone + '</p>' +
                    '<p class="text-sm"><strong>Количество человек:</strong> ' + data.people_count + '</p>' +
                    '<p class="text-sm"><strong>Дата и время:</strong> ' + slotDetails + '</p>' +
                    '<p class="text-sm mt-2"><strong>Статус:</strong> Ожидает подтверждения администратором в течение ' + holdMin + ' минут.</p>' +
                    '</div>';
                form.reset();
                selectedSlotData = null;
                loadSlots();
            })
            .catch(function(err) {
                var resultDiv = document.getElementById('booking-result');
                resultDiv.innerHTML = '<div class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded">' +
                    '<p class="font-semibold">Ошибка при создании бронирования</p>' +
                    '<p class="text-sm">' + (err.message || 'Пожалуйста, попробуйте еще раз') + '</p>' +
                    '</div>';
            })
            .finally(function() {
                btn.disabled = false;
                btn.textContent = 'Забронировать';
            });
        });
    }
});
