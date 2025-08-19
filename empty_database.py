import sqlite3

DB_PATH = 'service/database/app.db'

def empty_all_tables():
    conn = sqlite3.connect(DB_PATH)
    c = conn.cursor()
    # Disable foreign key checks for SQLite
    c.execute('PRAGMA foreign_keys = OFF;')
    # List of all tables to clear
    tables = ['messages', 'reactions', 'comments', 'contacts', 'conversations', 'users']
    for table in tables:
        c.execute(f'DELETE FROM {table}')
        print(f'Emptied {table}')
    conn.commit()
    c.execute('PRAGMA foreign_keys = ON;')
    conn.close()
    print('All tables emptied.')

if __name__ == '__main__':
    empty_all_tables()
