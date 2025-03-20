import websocket
import telebot
from telegram import InlineKeyboardButton, InlineKeyboardMarkup
from telegram.ext import Application, CommandHandler, CallbackQueryHandler
import threading
import time
import os
from dotenv import load_dotenv
import logging

load_dotenv()

logging.basicConfig(
    level=logging.INFO, 
    format="%(asctime)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

TOKEN = os.getenv("TELEGRAM_BOT_TOKEN")
WEBSOCKET_URL = os.getenv("WEBSOCKET_URL")

bot = telebot.TeleBot(token=TOKEN)

ws = None
chat_id = None

def on_message(ws, message):
    logger.info(f"Received from server: {message}")
    if chat_id:
        bot.send_message(chat_id=chat_id, text=message)
    else:
        logger.error("Chat ID is not set!")

def on_error(ws, error):
    logger.error(f"Error: {error}")

def on_close(ws, close_status_code, close_msg):
    logger.warning("Closed connection")
    connect_websocket()

def on_open(ws):
    logger.info("Connected to WebSocket server")

async def start(update, context):
    global chat_id
    chat_id = update.message.chat.id
    keyboard = [
        [InlineKeyboardButton("Узнать погоду", callback_data='weather')]
    ]
    reply_markup = InlineKeyboardMarkup(keyboard)
    await update.message.reply_text(
        'Добро пожаловать! Нажмите кнопку, чтобы узнать погоду в Санкт-Петербурге:', reply_markup=reply_markup
    )

# Функция для обработки нажатия на кнопку
async def button(update, context):
    query = update.callback_query
    await query.answer()

    # Проверяем, открыто ли WebSocket-соединение, если нет, переподключаемся
    if not ws or not ws.sock or not ws.sock.connected:
        logger.info("WebSocket connection is closed or not established. Reconnecting...")
        connect_websocket()

    # Отправляем запрос через WebSocket
    ws.send("weather")

def connect_websocket():
    global ws
    logger.info(f"Connecting to WebSocket server at {WEBSOCKET_URL}")
    ws = websocket.WebSocketApp(WEBSOCKET_URL,
                                 on_message=on_message,
                                 on_error=on_error,
                                 on_close=on_close,
                                 on_open=on_open)
    # Запускаем WebSocket в отдельном потоке
    threading.Thread(target=ws.run_forever).start()

def main():
    # Подключаем WebSocket
    connect_websocket()

    application = Application.builder().token(TOKEN).build()

    application.add_handler(CommandHandler('start', start))
    application.add_handler(CallbackQueryHandler(button))

    application.run_polling()

if __name__ == '__main__':
    main()
