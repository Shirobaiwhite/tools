print("Please enter an integer to check if that year is a leap year.")
year = input()
while not year.isnumeric():
    print('Please enter an integer.')
    year = input()
year = int(year)

def isLeapYear(year):
    if year < 1582:
        print('Please enter a year after 1581.')
        return -1
    if year > 9999:
        print('Please enter a year before 10000.')
        return -1
    return year % 400 == 0 or (year % 4 == 0 and year % 100 != 0)

if isLeapYear(year) == -1:
    year + 0
elif isLeapYear(year):
    print('You have entered a leap year.')
else:
    print('You have entered a non-leap year.')