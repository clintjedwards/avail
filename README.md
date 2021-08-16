# Avail(able)

This package allows the representation of a time frame using common cron syntax and
then efficiently checks whether a given golang time object exists within that time frame.

The package uses a subset of the extended cron standard:
https://en.wikipedia.org/wiki/Cron#CRON_expression

## Why

It is sometimes useful to represent certain timeframes using cron expressions.

A great example is an application that schedules employees for work. Representing the
case where an employee cannot work every year on their birthday, would be difficult without
representing it as some type of customized format (due to the nature that representing infinite
time would fill your database).

Using cron to achieve this allows the representation of situations like above to be compact and easy
to parse. Other advantages include that the cron format is well documented and can potentially
be represented in a user-friendly frontend component.

The drawback is that expressing things in cron is not always straight forward. For example, avoiding
scheduling an employee on their birthday every year is a single expression, but avoiding scheduling
an employee from Jun 1st to July 15th is two separate expressions.

## How

Avail implements/uses a stripped down version of the cron expression syntax as defined below.

    Field           Allowed values  Allowed special characters

    Minutes         0-59            * , -
    Hours           0-23            * , -
    Day of month    1-31            * , -
    Month           1-12            * , -
    Day of week     0-6             * , -
    Year            1970-2100       * , -

Avail accepts a cron expression in the format above, splits it into separate fields, parses it,
and generates map backed sets for each field in order to allow speedy checking of value existence.

This data structure allows avail to take a supplied time and check that each of the time's
elements exist in the representation of the cron expression.

## Usage

Initiate a new avail instance with cron expression.

    import "github.com/clintjedwards/avail/v2"

    avail, _ := avail.New("* * * * * *")

This will parse the cron expression given and
return a new `Avail` object containing your given expression and its parsed terms(each section
of the cron expression is called a term).

Then call `able` with a specified go time object.

    now := time.Now()

    fmt.Println(avail.Able(now))
    // Output: true
